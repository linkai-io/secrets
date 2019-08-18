# secrets
This is a demonstration package for accessing secured secrets from AWS SSM. Secrets should be created by an external service such as [binxio's secret-provider](https://github.com/binxio/cfn-secret-provider). Access to these secrets should be restricted by read only IAM policies. 

This package demonstrates how one could structure a secrets manager that allows for ease of development by substituting a secret store with environment variables (for running tests).

## Usage 
```Go
env := os.Getenv("APP_ENV")
region := os.Getenv("APP_REGION")

s := secrets.SecretsCache(env, region)
password, err := s.Password()
if err != nil {
  log.Fatalf("failed to get password: %v\n", err)
}
```

## Running the example locally:
- make buildexample
- docker run --name secrets -p 8000:8000 service:latest

## Running the example in AWS:
Build the docker container (make buildexample) and push to ECR. Ensure you have secrets configured via the [binxio secret-provider](https://github.com/binxio/cfn-secret-provider).  You'll need to create an ECS task (and cluster) with a policy similar to:

```YAML
  ServicePassword:
    Type: Custom::Secret
    Properties:
      Name: /some/${EnvironmentName}/password
      KeyAlias: alias/aws/ssm
      Alphabet: abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789
      Length: 30
      ReturnSecret: true
      ServiceToken: !Sub 'arn:aws:lambda:${AWS::Region}:${AWS::AccountId}:function:binxio-cfn-secret-provider'

  ServiceLogGroup:
    Type: AWS::Logs::LogGroup
    Properties:
      RetentionInDays: 30
      LogGroupName: !Sub "${EnvironmentName}-service"

  ServiceReadSSMRole:
    Type: AWS::IAM::Role
    Properties:
      AssumeRolePolicyDocument:
        Statement:
        - Effect: Allow
          Principal:
            Service: "ecs-tasks.amazonaws.com"
          Action: ['sts:AssumeRole']
      Path: /
      Policies:
      - PolicyName: !Sub "${EnvironmentName}-service-read-password"
        PolicyDocument:
          Statement:
          - Effect: Allow
            Action:
              - "ssm:Describe*"
              - "ssm:Get*"
              - "ssm:List*"
            Resource: !Sub "arn:aws:ssm:${AWS::Region}:${AWS::AccountId}:parameter/some/${EnvironmentName}/password"
  
  ServiceTaskDefinition:
    Type: AWS::ECS::TaskDefinition
    DependsOn: ["ServicePassword"]
    Properties:
      Family: service
      TaskRoleArn: !GetAtt 'ServiceReadSSMRole.Arn'
      NetworkMode: bridge
      ContainerDefinitions:
        - Name: service
          Image: FILL-ME-IN.dkr.ecr.us-east-1.amazonaws.com/service:latest
          Memory: 128
          Essential: true
          PortMappings:
            - ContainerPort: 8000
              Protocol: tcp
          Environment: 
            - Name: APP_ENV
              Value: !Sub "${EnvironmentName}"
            - Name: APP_REGION
              Value: !Ref "AWS::Region"
          LogConfiguration:
            LogDriver: 'awslogs'
            Options:
              awslogs-group: !Sub "${EnvironmentName}-service"
              awslogs-region: !Ref 'AWS::Region'
              awslogs-stream-prefix: 'service'
  
  ServiceDaemon:
    Type: AWS::ECS::Service
    Properties:
      ServiceName: 'service'
      Cluster: !Sub "${EnvironmentName}-cluster"
      TaskDefinition: !Ref 'ServiceTaskDefinition'
      SchedulingStrategy: 'DAEMON'
```

The above IAM policy restricts the ECS task to *only* being able to read the /some/environment/password SSM value. 