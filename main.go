package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "io/ioutil"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/ec2"
)

// Function to get the instance's public IP
func getPublicIP() (string, error) {
    resp, err := http.Get("http://169.254.169.254/latest/meta-data/public-ipv4")
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}

// Function to get the instance's private IP
func getPrivateIP() (string, error) {
    resp, err := http.Get("http://169.254.169.254/latest/meta-data/local-ipv4")
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}

// Function to get instance ID
func getInstanceID() (string, error) {
    resp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}

// Function to get attached volume IDs
func getAttachedVolumes(ec2Client *ec2.Client, instanceID string) ([]string, error) {
    input := &ec2.DescribeInstancesInput{
        InstanceIds: []string{instanceID},
    }

    result, err := ec2Client.DescribeInstances(context.TODO(), input)
    if err != nil {
        return nil, err
    }

    var volumeIDs []string
    for _, reservation := range result.Reservations {
        for _, instance := range reservation.Instances {
            for _, blockDevice := range instance.BlockDeviceMappings {
                volumeIDs = append(volumeIDs, aws.ToString(blockDevice.Ebs.VolumeId))
            }
        }
    }

    return volumeIDs, nil
}

func main() {
    // Get Public IP
    publicIP, err := getPublicIP()
    if err != nil {
        log.Fatalf("Error getting public IP: %v", err)
    }
    fmt.Printf("Public IP: %s\n", publicIP)

    // Get Private IP
    privateIP, err := getPrivateIP()
    if err != nil {
        log.Fatalf("Error getting private IP: %v", err)
    }
    fmt.Printf("Private IP: %s\n", privateIP)

    // Load the AWS config
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
    if err != nil {
        log.Fatalf("Unable to load SDK config, %v", err)
    }

    // Create EC2 client
    ec2Client := ec2.NewFromConfig(cfg)

    // Get Instance ID
    instanceID, err := getInstanceID()
    if err != nil {
        log.Fatalf("Error getting instance ID: %v", err)
    }
    fmt.Printf("Instance ID: %s\n", instanceID)

    // Get attached volume IDs
    volumeIDs, err := getAttachedVolumes(ec2Client, instanceID)
    if err != nil {
        log.Fatalf("Error getting volume IDs: %v", err)
    }
    fmt.Printf("Volume IDs: %v\n", volumeIDs)
}
