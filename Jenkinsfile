pipeline {
    agent any
    stages {
        stage('behave') {
            steps {
                sh 'make dockerbehave'
            }
        }
        stage('image') {
            steps {
                sh 'make image'
            }
        }
        stage('deploy') {
            steps {
                sh 'make deploy'
            }
        }
    }
}
