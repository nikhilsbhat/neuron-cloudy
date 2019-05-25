// Package loadbalancer contains methods that responds to and with tailor made input and outputs.
// To talk to cloud directly build wrapper around the respective interface as per the requirements.
package loadbalancer

import (
	"fmt"
	"strings"
	"time"

	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
)

// GetLoadbalancerInput implements the method GetAllLoadbalancer, GetAllClassicLb, Getloadbalancers, GetAllApplicationLb to fetch the granular level details of loadbalancers.
type GetLoadbalancerInput struct {

	//optional parameter; The names of the loadbalancers in array of which the information has to be fetched (both classic/network kind of loadbalancers).
	//one can omit this if he/she is passing ARN's of loadbalancers.
	//this parameter is mandatory if one wants to fetch the data of classic load balancers.
	LbNames []string `json:"lbnames,omitempty"`

	//optional parameter; The ARN's of the loadbalancers in array of which the information has to be fetched (only application kind of loadbalancers) one can omit this if he/she is passing names of loadbalancers.
	LbArns []string `json:"lbarns,omitempty"`

	//optional parameter if getallloadbalancer is used; Type of loadbalancers to fetch the appropriate data (classic/application).
	Type string `json:"Type,omitempty"`

	//optional parameter; Only when you need unfiltered result from cloud, enable this field by setting it to true. By default it is set to false.
	GetRaw bool `json:"getraw"`
}

// GetAllLoadbalancer will help in fetching information about all loadbalancers this include both applicaiton and classic.
func (lb *GetLoadbalancerInput) GetAllLoadbalancer(con aws.EstablishConnectionInput) ([]LoadBalanceResponse, error) {

	lbChn := make(chan interface{}, 2)
	defer close(lbChn)
	go func() {
		lbs, err := lb.GetAllClassicLb(con)
		time.Sleep(time.Second * 1)
		if err != nil {
			lbChn <- err
		} else {
			lbChn <- lbs
		}
	}()
	go func() {
		lbs, err := lb.GetAllApplicationLb(con)
		time.Sleep(time.Second * 2)
		if err != nil {
			lbChn <- err
		} else {
			lbChn <- lbs
		}
	}()

	// this is just a workaround and has to be fixed soon.
	classiclb := <-lbChn
	applicationlb := <-lbChn

	response := new(LoadBalanceResponse)
	switch calssic := classiclb.(type) {
	case []LoadBalanceResponse:
		response.ClassicLb = calssic
	case error:
		return nil, calssic
	default:
		return nil, fmt.Errorf("An unknown error occurred while returning classiclb data")
	}

	switch app := applicationlb.(type) {
	case []LoadBalanceResponse:
		response.ApplicationLb = app
	case error:
		return nil, app
	default:
		return nil, fmt.Errorf("An unknown error occurred while returning classiclb data")
	}

	resp := make([]LoadBalanceResponse, 0)
	resp = append(resp, *response)
	return resp, nil
}

// GetAllClassicLb will fetch information about all the classic loadbalancers.
func (lb *GetLoadbalancerInput) GetAllClassicLb(con aws.EstablishConnectionInput) ([]LoadBalanceResponse, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	// searching all the classic loadbalancer
	searchLbResult, searchErr := elb.DescribeAllClassicLoadbalancer(
		&aws.DescribeLoadbalancersInput{},
	)
	if searchErr != nil {
		return nil, searchErr
	}

	lbList := make([]LoadBalanceResponse, 0)

	if lb.GetRaw == true {
		lbList = append(lbList, LoadBalanceResponse{GetClassicLbsRaw: searchLbResult})
		return lbList, nil
	}
	for _, load := range searchLbResult.LoadBalancerDescriptions {
		lbList = append(lbList, LoadBalanceResponse{Name: *load.LoadBalancerName, LbDns: *load.DNSName, Createdon: (*load.CreatedTime).String(), Type: "classic", Scheme: *load.Scheme, VpcId: *load.VPCId})
	}
	return lbList, nil
}

// GetAllApplicationLb will help in fetching the information about all applicaiton loadbalancer present in the region selected.
func (lb *GetLoadbalancerInput) GetAllApplicationLb(con aws.EstablishConnectionInput) ([]LoadBalanceResponse, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}
	// searching all the application loadbalancer
	searchLbResult, searchErr := elb.DescribeAllApplicationLoadbalancer(
		&aws.DescribeLoadbalancersInput{},
	)
	if searchErr != nil {
		return nil, searchErr
	}

	lbList := make([]LoadBalanceResponse, 0)

	for _, load := range searchLbResult.LoadBalancers {

		// searching target group for the corresponding loadbalancer
		searchTarget, tarerr := elb.DescribeTargetgroups(
			&aws.DescribeLoadbalancersInput{
				LbArns: []string{*load.LoadBalancerArn},
			},
		)
		if tarerr != nil {
			return nil, tarerr
		}

		// searching listeners for the corresponding loadbalancer
		searchListeners, lisserr := elb.DescribeListners(
			&aws.DescribeLoadbalancersInput{
				LbArns: []string{*load.LoadBalancerArn},
			},
		)
		if lisserr != nil {
			return nil, lisserr
		}

		response := new(LoadBalanceResponse)
		if lb.GetRaw == true {
			response.GetApplicationLbRaw.GetApplicationLbRaw = load
			response.GetApplicationLbRaw.GetTargetGroupRaw = searchTarget
			response.GetApplicationLbRaw.GetListnersRaw = searchListeners
			lbList = append(lbList, *response)
		} else {

			tarArn := make([]string, 0)
			for _, tar := range searchTarget.TargetGroups {
				tarArn = append(tarArn, *tar.TargetGroupArn)
			}

			lisRrn := make([]string, 0)
			for _, lis := range searchListeners.Listeners {
				lisRrn = append(lisRrn, *lis.ListenerArn)
			}

			response.Name = *load.LoadBalancerName
			response.LbDns = *load.DNSName
			response.LbArn = *load.LoadBalancerArn
			response.Createdon = (*load.CreatedTime).String()
			response.Type = *load.Type
			response.Scheme = *load.Scheme
			response.VpcId = *load.VpcId
			response.TargetArn = tarArn
			response.ListnerArn = lisRrn
			lbList = append(lbList, *response)
		}
	}
	return lbList, nil
}

// Getloadbalancers will fetch information about the type loadbalancer selected.
func (lb *GetLoadbalancerInput) Getloadbalancers(con aws.EstablishConnectionInput) ([]LoadBalanceResponse, error) {

	switch strings.ToLower(lb.Type) {
	case "classic":
		getlb, err := lb.GetClassicloadbalancers(con)
		if err != nil {
			return nil, err
		}
		return getlb, nil
	case "application":
		getlb, err := lb.GetApplicationloadbalancers(con)
		if err != nil {
			return nil, err
		}
		return getlb, nil
	default:
		return nil, fmt.Errorf("You provided unknown loadbalancer type, enter a valid LB type")
	}
}

// GetClassicloadbalancers will help in fetching the information of the selected classic loadbalancer.
func (lb *GetLoadbalancerInput) GetClassicloadbalancers(con aws.EstablishConnectionInput) ([]LoadBalanceResponse, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	lbresponse, err := elb.DescribeClassicLoadbalancer(
		&aws.DescribeLoadbalancersInput{
			LbNames: lb.LbNames,
		},
	)
	if err != nil {
		return nil, err
	}
	lbResponse := make([]LoadBalanceResponse, 0)

	if lb.GetRaw == true {
		lbResponse = append(lbResponse, LoadBalanceResponse{GetClassicLbsRaw: lbresponse})
	}

	for _, lb := range lbresponse.LoadBalancerDescriptions {
		response := new(LoadBalanceResponse)
		response.Name = *load.LoadBalancerName
		response.LbDns = *load.DNSName
		response.Createdon = (*load.CreatedTime).String()
		response.Type = "classic"
		response.Scheme = *load.Scheme
		response.VpcId = *load.VPCId
		lbResponse = append(lbResponse, *response)
	}

	return lbResponse, nil
}

// GetApplicationloadbalancers will help in fetching the information of the selected application load balancer.
func (lb *GetLoadbalancerInput) GetApplicationloadbalancers(con aws.EstablishConnectionInput) ([]LoadBalanceResponse, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	lbin := new(aws.DescribeLoadbalancersInput)
	lbin.LbNames = lb.LbNames
	lbin.LbArns = lb.LbArns
	getLoadbalancer, err := elb.DescribeApplicationLoadbalancer(lbin)
	if err != nil {
		return nil, err
	}
	lbResponse := make([]LoadBalanceResponse, 0)
	for _, load := range getLoadbalancer.LoadBalancers {

		// searching target group for the corresponding loadbalancer
		lbin.LbArns = []string{*load.LoadBalancerArn}
		searchTarget, tarerr := elb.DescribeTargetgroups(lbin)
		if tarerr != nil {
			return nil, tarerr
		}

		// searching listeners for the corresponding loadbalancer
		searchListeners, lisserr := elb.DescribeListners(lbin)
		if lisserr != nil {
			return nil, lisserr
		}

		response := new(LoadBalanceResponse)
		if lb.GetRaw == true {
			response.GetApplicationLbRaw.GetApplicationLbRaw = load
			response.GetApplicationLbRaw.GetTargetGroupRaw = searchTarget
			response.GetApplicationLbRaw.GetListnersRaw = searchListeners
			lbResponse = append(lbResponse, *response)
		} else {
			tarArn := make([]string, 0)
			for _, tar := range searchTarget.TargetGroups {
				tarArn = append(tarArn, *tar.TargetGroupArn)
			}

			lisRrn := make([]string, 0)
			for _, lis := range searchListeners.Listeners {
				lisRrn = append(lisRrn, *lis.ListenerArn)
			}

			response.Name = *load.LoadBalancerName
			response.LbDns = *load.DNSName
			response.LbArn = *load.LoadBalancerArn
			response.Createdon = (*load.CreatedTime).String()
			response.Type = *load.Type
			response.Scheme = *load.Scheme
			response.VpcId = *load.VpcId
			response.TargetArn = tarArn
			response.ListnerArn = lisRrn
			lbResponse = append(lbResponse, *response)
		}
	}
	return lbResponse, nil
}

// FindClassicLoadbalancer will help in fetching the information about the selested classic loadbalancer.
func (lb *GetLoadbalancerInput) FindClassicLoadbalancer(con aws.EstablishConnectionInput) (bool, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return false, sesserr
	}

	getLoadbalancer, err := elb.DescribeClassicLoadbalancer(
		&aws.DescribeLoadbalancersInput{
			LbNames: lb.LbNames,
		},
	)
	if err != nil {
		return false, err
	}
	if len(getLoadbalancer.LoadBalancerDescriptions) != 0 {
		return true, nil
	}
	return false, fmt.Errorf("Could not find the entered loadbalancer, please enter valid/existing loadbalancer Name")
}

// FindApplicationLoadbalancer will return ture if loadbalancer exists in the system.
func (lb *GetLoadbalancerInput) FindApplicationLoadbalancer(con aws.EstablishConnectionInput) (bool, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return false, sesserr
	}

	getLoadbalancer, err := elb.DescribeApplicationLoadbalancer(
		&aws.DescribeLoadbalancersInput{
			LbNames: lb.LbNames,
			LbArns:  lb.LbArns,
		},
	)
	if err != nil {
		return false, err
	}
	if len(getLoadbalancer.LoadBalancers) != 0 {
		return true, nil
	}
	return false, fmt.Errorf("Could not find the entered loadbalancer, please enter valid/existing loadbalancer ARN/Name")
}

// GetArnFromLoadbalancer will help in fetching ARN from the selected loadbalancer.
func (lb *GetLoadbalancerInput) GetArnFromLoadbalancer(con aws.EstablishConnectionInput) (LoadBalanceResponse, error) {

	//get the relative sessions before proceeding further
	elb, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return LoadBalanceResponse{}, sesserr
	}

	getLoadbalancer, err := elb.DescribeApplicationLoadbalancer(
		&aws.DescribeLoadbalancersInput{
			LbNames: lb.LbNames,
		},
	)
	if err != nil {
		return LoadBalanceResponse{}, err
	}
	arns := make([]string, 0)
	for _, lb := range getLoadbalancer.LoadBalancers {
		arns = append(arns, *lb.LoadBalancerArn)
	}

	return LoadBalanceResponse{LbArns: arns}, nil
}
