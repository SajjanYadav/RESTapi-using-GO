pipeline {
    agent any
    tools{
        go '1.20.4'
    }

    environment {
        GO111MODULE='on'
    }

    stages{
        stage('Test') {
            steps{
                git 'https://github.com/SajjanYadav/RESTapi-using-GO.git'    
                sh 'go test ./...'                                           
            }
        }
    }
}