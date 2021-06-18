package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"net"
	"regexp"
)

var validateName = validation.StringMatch(
	regexp.MustCompile(`^[A-Za-z0-9\p{Han}-_]{1,63}$`),
	"expected value to be 1 - 63 characters and only support chinese, english, numbers, '-', '_'",
)

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
	}

	return
}

func validateIpAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	res := net.ParseIP(value)

	if res == nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid ip address, got error parsing: %s", k, value))
	}

	return
}

func validateSubnetType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "Reserve" && value != "Normal" && value != "Physical" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid subnet type, got error parsing: %s", k, value))
	}
	return
}

func validateLbState(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "start" && value != "stop" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid loadbalancer state, got error parsing: %s", k, value))
	}
	return
}
func validateLbType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "public" && value != "internal" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid loadbalancer type, got error parsing: %s", k, value))
	}
	return
}

func validateRouteType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "InternetGateway" && value != "Tunnel" && value != "Host" && value != "Peering" && value != "DirectConnect" && value != "Vpn" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid route type, got error parsing: %s", k, value))
	}
	return
}

func validateNatType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "public" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid nat type, got error parsing: %s", k, value))
	}
	return
}

func validateNatMode(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "Vpc" && value != "Subnet" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid nat mode, got error parsing: %s", k, value))
	}
	return
}

func validateNatIpNumber(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 10 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid nat ip number in 1-10 and control by quota system, got error parsing: %d", k, value))
	}
	return
}

func validateNatBandWidth(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 15000 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid nat bandwidth in 1-15000 and control by quota system, got error parsing: %d", k, value))
	}
	return
}

func validateNatChargeType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "Peak" && value != "Daily" && value != "TrafficMonthly" &&
		value != "DailyPaidByTransfer" && value != "PostPaidByAdvanced95Peak" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid nat charge type and control by price system, got error parsing: %s", k, value))
	}
	return
}

func validateKecSystemDiskType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "Local_SSD" && value != "SSD3.0" && value != "EHDD" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid System Disk Type and control by price system, got error parsing: %s", k, value))
	}
	return
}

func validateKecSystemDiskSize(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 500 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid System Disk Size and control by price system, got error parsing: %d", k, value))
	}
	return
}

func validateKecDataDiskType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "SSD3.0" && value != "EHDD" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid Data Disk Type and control by price system, got error parsing: %s", k, value))
	}
	return
}

func validateKecDataDiskSize(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 10 || value > 16000 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid Data Disk Size and control by price system, got error parsing: %d", k, value))
	}
	return
}

func validateKecInstanceAgent(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value != 0 && value != 1 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid Instance Agent and control by price system, got error parsing: %d", k, value))
	}
	return
}

func validateKecScalingGroupSize(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 1000 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid ScalingGroup min or max size, got error parsing: %d", k, value))
	}
	return
}

func validateKecScalingGroupDesiredCapacity(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 1000 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid ScalingGroupDesiredCapacity, got error parsing: %d", k, value))
	}
	return
}

func validateKecScalingGroupRemovePolicy(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "RemoveOldestInstance" && value != "RemoveNewestInstance" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid ScalingGroupRemovePolicy, got error parsing: %s", k, value))
	}
	return
}

func validateKecScalingGroupSubnetStrategy(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "balanced-distribution" && value != "choice-first" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid ScalingGroupSubnetStrategy, got error parsing: %s", k, value))
	}
	return
}

func validateKecScalingGroupStatus(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "UnActive" && value != "Active" {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid ScalingGroupStatus, got error parsing: %s", k, value))
	}
	return
}

func validatePurchaseTime(req *map[string]interface{}, purchaseTimeField string, chargeTypeField string, chargeTypes []string) error {
	if v, ok := (*req)[chargeTypeField]; ok {
		flag := false
		for _, t := range chargeTypes {
			if t == v {
				flag = true
				if _, ok := (*req)[purchaseTimeField]; !ok {
					return fmt.Errorf(
						"%q must contain a value", purchaseTimeField)
				}
			}
		}
		if _, ok := (*req)[purchaseTimeField]; ok {
			if !flag {
				delete(*req, purchaseTimeField)
			}
		}
	}
	return nil
}

//校验Ks3 Bucket name
/*
func validateKs3BucketName(value string) error {
	if (len(value) < 3) || (len(value) > 63) { //3~63字符之间
		return fmt.Errorf("%q must contain from 3 to 63 characters", value)
	}
	if !regexp.MustCompile(`^[0-9a-z-.]+$`).MatchString(value) { //小写和数字
		return fmt.Errorf("only lowercase alphanumeric characters and hyphens allowed in %q", value)
	}
	if regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`).MatchString(value) { //不能是IP
		return fmt.Errorf("%q must not be formatted as an IP address", value)
	}
	if strings.HasPrefix(value, `.`) { //不能以点开头
		return fmt.Errorf("%q cannot start with a period", value)
	}
	if strings.HasSuffix(value, `.`) { //不能以点结尾
		return fmt.Errorf("%q cannot end with a period", value)
	}
	if strings.Contains(value, `..`) { //不能包含两个点
		return fmt.Errorf("%q can be only one period between labels", value)
	}
	return nil
}

func validateKs3BucketLifecycleTransitionStorageClass() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		s3.TransitionStorageClassGlacier,
		s3.TransitionStorageClassStandardIa,
		s3.TransitionStorageClassOnezoneIa,
		s3.TransitionStorageClassIntelligentTiering,
		s3.TransitionStorageClassDeepArchive,
	}, false)
}
func validateKs3BucketLifecycleTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", value))
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as RFC3339 Timestamp Format", value))
	}

	return
}

*/
