package loadbalancer

import (
	"fmt"
	"strings"
	"time"

	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
)

// DeleteLoadbalancerInput will values enough to delete a loadbalancer
type DeleteLoadbalancerInput struct {

	//optional parameter; The name of the loadbalancers which has to be deleted (only classic/network kind of loadbalancers) one can omit this if he/she needs to delete application loadbalncers.
	LbNames []string `json:"LbDns,omitempty"`

	//optional parameter; The ARN's of the loadbalancers which has to be deleted (only application kind of loadbalancers) one can omit this if he/she needs to delete classic/network loadbalncers.
	LbArns []string `json:"LbArn,omitempty"`

	//mandatory parameter; Type of loadbalancers to delete the appropriate one (classic/application).
	Type string `json:"LbArns,omitempty"`

	//optional parameter; Only when you need unfiltered result from cloud, enable this field by setting it to true. By default it is set to false.
	GetRaw bool `json:"GetRaw"`
}

// LoadBalanceDeleteResponse contains filtered/unfiltered response obtained from DeleteLoadbalancer
type LoadBalanceDeleteResponse struct {
	LbDeleteStatus string `json:"LbDeleteStatus,omitempty"`
	LbArn          string `json:"LbArn,omitempty"`
	LbName         string `json:"LbName,omitempty"`
}

// DeleteLoadbalancer is actually responsible for deleting loadbalancer asper ther details passed to it.
func (d *DeleteLoadbalancerInput) DeleteLoadbalancer(con aws.EstablishConnectionInput) ([]LoadBalanceDeleteResponse, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	lbDeleteStatus := make([]LoadBalanceDeleteResponse, 0)
	switch strings.ToLower(d.Type) {
	case "application":

		if (d.LbNames != nil) && (d.LbArns != nil) {
			return nil, fmt.Errorf("You provided both LbNames and LbArns to fetch applicationlb data, has to provide either of them")
		}

		lbin := new(GetLoadbalancerInput)
		lbin.LbNames = d.LbNames
		lbin.LbArns = d.LbArns
		_, err := lbin.FindApplicationLoadbalancer(con)
		if err != nil {
			return nil, err
		}

		if d.LbNames != nil {
			for _, lb := range d.LbNames {
				//fetching arn of loadbalancer from name
				lbarnin := GetLoadbalancerInput{LbNames: []string{lb}}
				lbarn, arnerr := lbarnin.GetArnFromLoadbalancer(con)
				if arnerr != nil {
					return nil, arnerr
				}

				getlb := new(aws.DescribeLoadbalancersInput)
				getlb.LbArns = lbarn.LbArns

				//fetching arn of targetgroup
				tararn, tararnerr := elb.DescribeTargetgroups(getlb)
				if tararnerr != nil {
					return nil, tararnerr
				}

				//fetching arn of listeners
				lisarn, lisarnerr := elb.DescribeListners(getlb)
				if lisarnerr != nil {
					return nil, lisarnerr
				}

				delb := new(aws.DeleteLoadbalancerInput)
				if lisarn.Listeners != nil {
					delb.ListenerArn = *lisarn.Listeners[0].ListenerArn
					//deletion of listeners
					deliserr := elb.DeleteAppListeners(delb)
					if deliserr != nil {
						return nil, deliserr
					}
				}

				//deleting loadbalancers
				delb.LbArn = lbarn.LbArns[0]
				delerr := elb.DeleteAppLoadbalancer(delb)
				if delerr != nil {
					return nil, delerr
				}

				//waiting till the loadbalancer gets deleted completed
				waiterr := elb.WaitTillLbDeletionSuccessfull(getlb)
				if waiterr != nil {
					return nil, waiterr
				}

				// making this to sleep is not good idea but this is temporary fix
				time.Sleep(5 * time.Second)
				//deletion of targetgroups
				delb.TargetArn = *tararn.TargetGroups[0].TargetGroupArn
				tarerr := elb.DeleteTargetGroup(delb)
				if tarerr != nil {
					return nil, tarerr
				}

				lbDeleteStatus = append(lbDeleteStatus, LoadBalanceDeleteResponse{LbDeleteStatus: "LoadBalancer deletion is successful", LbArn: lbarn.LbArns[0]})
			}
			return lbDeleteStatus, nil
		}

		if d.LbArns != nil {

			for _, lbarn := range d.LbArns {
				//deleting loadbalancers
				delb := new(aws.DeleteLoadbalancerInput)
				delb.LbArn = lbarn
				delerr := elb.DeleteAppLoadbalancer(delb)
				if delerr != nil {
					return nil, delerr
				}

				deslb := new(aws.DescribeLoadbalancersInput)
				deslb.LbArns = []string{lbarn}
				//waiting till the loadbalancer gets deleted completed
				waiterr := elb.WaitTillLbDeletionSuccessfull(deslb)
				if waiterr != nil {
					return nil, waiterr
				}

				//fetching arn of targetgroup
				tararn, tararnerr := elb.DescribeTargetgroups(deslb)
				if tararnerr != nil {
					return nil, tararnerr
				}

				//deletion of targetgroups
				delb.TargetArn = *tararn.TargetGroups[0].TargetGroupArn
				tarerr := elb.DeleteTargetGroup(delb)
				if tarerr != nil {
					return nil, tarerr
				}

				lbDeleteStatus = append(lbDeleteStatus, LoadBalanceDeleteResponse{LbDeleteStatus: "LoadBalancer deletion is successful", LbArn: lbarn})
			}
			return lbDeleteStatus, nil
		}
		return nil, fmt.Errorf("You selected application lb to fetch data, but I couldn't find any valid inputs. Its is empty input")
	case "classic":
		if d.LbNames != nil {

			lbin := GetLoadbalancerInput{LbNames: d.LbNames}
			_, err := lbin.FindClassicLoadbalancer(con)
			if err != nil {
				return nil, err
			}
			for _, lbname := range d.LbNames {
				delErr := elb.DeleteClassicLoadbalancer(
					&aws.DeleteLoadbalancerInput{
						LbName: lbname,
					},
				)
				if delErr != nil {
					return nil, delErr
				}

				lbDeleteStatus = append(lbDeleteStatus, LoadBalanceDeleteResponse{LbDeleteStatus: "LoadBalancer deletion is successful", LbName: lbname})
			}
			return lbDeleteStatus, nil
		}
		return nil, fmt.Errorf("You selected classic lb to fetch data, but I couldn't find any valid inputs. Its is empty input")
	default:
		return nil, fmt.Errorf("You provided empty struct to delete loadbalancer, you have to pass either Names or Arns")
	}
}
