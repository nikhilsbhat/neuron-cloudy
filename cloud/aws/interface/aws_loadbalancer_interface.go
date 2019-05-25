package neuronaws

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	err "github.com/nikhilsbhat/neuron-cloudy/error"
)

// LoadBalanceCreateInput implements various methods to create variours types of loadbalancers in aws.
// It holds all the required values for the same.
type LoadBalanceCreateInput struct {
	Name              string
	VpcId             string
	Subnets           []string
	AvailabilityZones []string
	SecurityGroups    []string
	Scheme            string
	Type              string
	SslCert           string
	SslPolicy         string
	LbPort            int64
	InstPort          int64
	Lbproto           string
	Instproto         string
	HttpCode          string
	HealthPath        string
	IpAddressType     string
	TargetArn         string
	LbArn             string
}

// LoadBalanceResponse returns the filtered/unfiltered results obtained from aws.
type LoadBalanceResponse struct {
	Name            string                `json:"name,omitempty"`
	Type            string                `json:"type,omitempty"`
	LbDns           string                `json:"lbdns,omitempty"`
	LbArn           string                `json:"lbarn,omitempty"`
	TargetArn       string                `json:"targetarn,omitempty"`
	ListnerArn      string                `json:"listnerarn,omitempty"`
	Createdon       string                `json:"createdon,omitempty"`
	VpcId           string                `json:"vpcid,omitempty"`
	Scheme          string                `json:"scheme,omitempty"`
	DefaultResponse interface{}           `json:"defaultresponse,omitempty"`
	LbDeleteStatus  string                `json:"lbdeletestatus,omitempty"`
	ApplicationLb   []LoadBalanceResponse `json:"applicationlb,omitempty"`
	ClassicLb       []LoadBalanceResponse `json:"classiclb,omitempty"`
}

// DeleteLoadbalancerInput implements various methods to delete various types of load balancers.
type DeleteLoadbalancerInput struct {
	LbName      string
	LbArn       string
	TargetArn   string
	ListenerArn string
}

// DescribeLoadbalancersInput implements various methods to fetch details of various types of load balancers.
type DescribeLoadbalancersInput struct {
	LbNames     []string
	LbArns      []string
	TargetArns  []string
	ListnerArns []string
}

// CreateClassicLb helps in creating load balancer of type classic
func (sess *EstablishedSession) CreateClassicLb(lb LoadBalanceCreateInput) (*elb.CreateLoadBalancerOutput, error) {

	if sess.Elb != nil {
		listeners := make([]*elb.Listener, 0)
		switch lb.Lbproto {
		case "HTTP":
			listeners = append(listeners, &elb.Listener{
				InstancePort:     aws.Int64(lb.InstPort),
				InstanceProtocol: aws.String(lb.Instproto),
				LoadBalancerPort: aws.Int64(lb.LbPort),
				Protocol:         aws.String(lb.Lbproto),
			})
		case "HTTPS":
			listeners = append(listeners, &elb.Listener{
				InstancePort:     aws.Int64(lb.InstPort),
				InstanceProtocol: aws.String(lb.Instproto),
				LoadBalancerPort: aws.Int64(lb.LbPort),
				Protocol:         aws.String(lb.Lbproto),
				SSLCertificateId: aws.String(lb.SslCert),
			})
		default:
			return nil, fmt.Errorf("You provided unknown loadbalancer protocol, enter a valid protocol")
		}
		input := &elb.CreateLoadBalancerInput{
			Listeners:        listeners,
			LoadBalancerName: aws.String(lb.Name),
			Scheme:           aws.String(lb.Scheme),
			SecurityGroups:   aws.StringSlice(lb.SecurityGroups),
			Subnets:          aws.StringSlice(lb.Subnets),
		}
		result, err := (sess.Elb).CreateLoadBalancer(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

// CreateApplicationLb helps in creating loadbalancer of type application
func (sess *EstablishedSession) CreateApplicationLb(lb LoadBalanceCreateInput) (*elbv2.CreateLoadBalancerOutput, error) {

	if sess.Elb2 != nil {
		input := &elbv2.CreateLoadBalancerInput{
			Name:           aws.String(lb.Name),
			Scheme:         aws.String(lb.Scheme),
			Subnets:        aws.StringSlice(lb.Subnets),
			SecurityGroups: aws.StringSlice(lb.SecurityGroups),
			IpAddressType:  aws.String(lb.IpAddressType),
			Tags: []*elbv2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(lb.Name),
				}},
		}

		result, err := (sess.Elb2).CreateLoadBalancer(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()

}

// CreateTargetGroups helps in creating target group, the group which will be attached to the loadbalancer that will be created.
// These target groups consists of VM's which actually take the load.
func (sess *EstablishedSession) CreateTargetGroups(lb *LoadBalanceCreateInput) (*elbv2.CreateTargetGroupOutput, error) {

	if sess.Elb2 != nil {
		input := &elbv2.CreateTargetGroupInput{
			Name:                       aws.String(lb.Name),
			Port:                       aws.Int64(lb.LbPort),
			Protocol:                   aws.String(lb.Lbproto),
			VpcId:                      aws.String(lb.VpcId),
			HealthCheckProtocol:        aws.String(lb.Instproto),
			HealthCheckPort:            aws.String(strconv.FormatInt(lb.InstPort, 10)),
			HealthCheckPath:            aws.String(lb.HealthPath),
			HealthCheckIntervalSeconds: aws.Int64(30),
			HealthCheckTimeoutSeconds:  aws.Int64(5),
			HealthyThresholdCount:      aws.Int64(5),
			UnhealthyThresholdCount:    aws.Int64(2),
			Matcher:                    &elbv2.Matcher{HttpCode: &lb.HttpCode},
		}

		result, err := (sess.Elb2).CreateTargetGroup(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

// CreateApplicationListners helps in creating listners for the loadbalancers.
func (sess *EstablishedSession) CreateApplicationListners(lb *LoadBalanceCreateInput) (*elbv2.CreateListenerOutput, error) {

	if sess.Elb2 != nil {
		var input elbv2.CreateListenerInput
		switch lb.Lbproto {
		case "HTTP":
			input = elbv2.CreateListenerInput{
				DefaultActions: []*elbv2.Action{
					{
						TargetGroupArn: aws.String(lb.TargetArn),
						Type:           aws.String("forward"),
					},
				},
				LoadBalancerArn: aws.String(lb.LbArn),
				Port:            aws.Int64(lb.LbPort),
				Protocol:        aws.String(lb.Lbproto),
			}
		case "HTTPS":
			input = elbv2.CreateListenerInput{
				Certificates: []*elbv2.Certificate{
					{
						CertificateArn: aws.String(lb.SslCert),
					},
				},
				DefaultActions: []*elbv2.Action{
					{
						TargetGroupArn: aws.String(lb.TargetArn),
						Type:           aws.String("forward"),
					},
				},
				LoadBalancerArn: aws.String(lb.LbArn),
				Port:            aws.Int64(lb.LbPort),
				Protocol:        aws.String(lb.Lbproto),
				SslPolicy:       aws.String(lb.SslPolicy),
			}
		default:
			return nil, fmt.Errorf("You provided unknown loadbalancer protocol, enter a valid protocol")
		}

		result, err := (sess.Elb2).CreateListener(&input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

// DescribeClassicLoadbalancer describes the details of the selected classic loadbalancers.
func (sess *EstablishedSession) DescribeClassicLoadbalancer(lb *DescribeLoadbalancersInput) (*elb.DescribeLoadBalancersOutput, error) {

	if sess.Elb != nil {
		if lb.LbNames != nil {
			input := &elb.DescribeLoadBalancersInput{
				LoadBalancerNames: aws.StringSlice(lb.LbNames),
			}
			result, err := (sess.Elb).DescribeLoadBalancers(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf(fmt.Sprintf("%v DescribeClassicLoadbalancer", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()
}

// DescribeAllClassicLoadbalancer describes the details of all the classic loadbalancers in the selected region.
func (sess *EstablishedSession) DescribeAllClassicLoadbalancer(lb *DescribeLoadbalancersInput) (*elb.DescribeLoadBalancersOutput, error) {

	if sess.Elb != nil {
		input := &elb.DescribeLoadBalancersInput{}
		result, err := (sess.Elb).DescribeLoadBalancers(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

// DescribeApplicationLoadbalancer describes the details of the selected application loadbalancers.
func (sess *EstablishedSession) DescribeApplicationLoadbalancer(lb *DescribeLoadbalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {

	if sess.Elb2 != nil {
		if lb.LbArns != nil {
			input := &elbv2.DescribeLoadBalancersInput{
				LoadBalancerArns: aws.StringSlice(lb.LbArns),
			}
			result, err := (sess.Elb2).DescribeLoadBalancers(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if lb.LbNames != nil {
			input := &elbv2.DescribeLoadBalancersInput{
				Names: aws.StringSlice(lb.LbNames),
			}
			result, err := (sess.Elb2).DescribeLoadBalancers(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}

		return nil, fmt.Errorf(fmt.Sprintf("%v DescribeApplicationLoadbalancer", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()
}

// DescribeAllApplicationLoadbalancer describes the details of all the application loadbalancers in the selected region.
func (sess *EstablishedSession) DescribeAllApplicationLoadbalancer(lb *DescribeLoadbalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {

	if sess.Elb2 != nil {
		input := &elbv2.DescribeLoadBalancersInput{}
		result, err := (sess.Elb2).DescribeLoadBalancers(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

// DescribeTargetgroups describes the target group of the loadbalancers.
func (sess *EstablishedSession) DescribeTargetgroups(lb *DescribeLoadbalancersInput) (*elbv2.DescribeTargetGroupsOutput, error) {
	if sess.Elb2 != nil {
		if lb.TargetArns != nil {
			input := &elbv2.DescribeTargetGroupsInput{
				TargetGroupArns: aws.StringSlice(lb.TargetArns),
			}
			result, err := (sess.Elb2).DescribeTargetGroups(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if lb.LbArns != nil {
			input := &elbv2.DescribeTargetGroupsInput{
				LoadBalancerArn: aws.String(lb.LbArns[0]),
			}
			result, err := (sess.Elb2).DescribeTargetGroups(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf(fmt.Sprintf("%v DescribeTargetgroups", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()
}

// DescribeAllTargetgroups describes all the target group present in the selected region.
func (sess *EstablishedSession) DescribeAllTargetgroups(lb *DescribeLoadbalancersInput) (*elbv2.DescribeTargetGroupsOutput, error) {

	if sess.Elb2 != nil {
		input := &elbv2.DescribeTargetGroupsInput{}
		result, err := (sess.Elb2).DescribeTargetGroups(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()
}

// DescribeListners helps in describing the selected listners.
func (sess *EstablishedSession) DescribeListners(lb *DescribeLoadbalancersInput) (*elbv2.DescribeListenersOutput, error) {

	if sess.Elb2 != nil {
		if lb.ListnerArns != nil {
			input := &elbv2.DescribeListenersInput{
				ListenerArns: aws.StringSlice(lb.ListnerArns),
			}
			result, err := (sess.Elb2).DescribeListeners(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}

		if lb.LbArns != nil {
			input := &elbv2.DescribeListenersInput{
				LoadBalancerArn: aws.String(lb.LbArns[0]),
			}
			result, err := (sess.Elb2).DescribeListeners(input)

			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, fmt.Errorf(fmt.Sprintf("%v DescribeListners", err.EmptyStructError()))
	}
	return nil, err.InvalidSession()

}

// DescribeAllListners helps in describing all the listners in the selscted region in aws.
func (sess *EstablishedSession) DescribeAllListners(lb *DescribeLoadbalancersInput) (*elbv2.DescribeListenersOutput, error) {

	if sess.Elb2 != nil {
		input := &elbv2.DescribeListenersInput{}
		result, err := (sess.Elb2).DescribeListeners(input)

		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, err.InvalidSession()

}

// DeleteClassicLoadbalancer will delete the selected loadbalancer of type classic in cloud aws.
func (sess *EstablishedSession) DeleteClassicLoadbalancer(lb *DeleteLoadbalancerInput) error {

	if sess.Elb != nil {
		input := &elb.DeleteLoadBalancerInput{
			LoadBalancerName: aws.String(lb.LbName),
		}
		_, err := (sess.Elb).DeleteLoadBalancer(input)

		if err != nil {
			return err
		}
		return nil
	}
	return err.InvalidSession()

}

// DeleteAppLoadbalancer will delete the selected loadbalancer of type application in cloud aws.
func (sess *EstablishedSession) DeleteAppLoadbalancer(lb *DeleteLoadbalancerInput) error {

	if sess.Elb2 != nil {
		input := &elbv2.DeleteLoadBalancerInput{
			LoadBalancerArn: aws.String(lb.LbArn),
		}
		_, err := (sess.Elb2).DeleteLoadBalancer(input)

		if err != nil {
			return err
		}
		return nil
	}
	return err.InvalidSession()
}

// DeleteTargetGroup deletes the selected target group in the cloud aws.
func (sess *EstablishedSession) DeleteTargetGroup(lb *DeleteLoadbalancerInput) error {

	if sess.Elb2 != nil {
		if lb.TargetArn != "" {
			input := &elbv2.DeleteTargetGroupInput{
				TargetGroupArn: aws.String(lb.TargetArn),
			}
			_, err := (sess.Elb2).DeleteTargetGroup(input)

			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v DeleteTargetGroup", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

// DeleteAppListeners helps in deleting the selected application listners of the loadbalancers.
func (sess *EstablishedSession) DeleteAppListeners(lb *DeleteLoadbalancerInput) error {

	if sess.Elb2 != nil {
		if lb.ListenerArn != "" {
			input := &elbv2.DeleteListenerInput{
				ListenerArn: aws.String(lb.ListenerArn),
			}
			_, err := (sess.Elb2).DeleteListener(input)

			if err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("%v DeleteAppListeners", err.EmptyStructError()))
	}
	return err.InvalidSession()
}

// WaitTillLbDeletionSuccessfull will make the method who called this to wait till the deletion of loadbalancer is successfull.
func (sess *EstablishedSession) WaitTillLbDeletionSuccessfull(lb *DescribeLoadbalancersInput) error {

	if sess.Elb2 != nil {
		input := &elbv2.DescribeLoadBalancersInput{
			LoadBalancerArns: aws.StringSlice(lb.LbArns),
		}
		err := (sess.Elb2).WaitUntilLoadBalancersDeleted(input)

		if err != nil {
			return err
		}
		return nil
	}
	return err.InvalidSession()
}
