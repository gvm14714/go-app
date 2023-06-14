pipeline {
    agent any
    environment {
        DOCKERHUB_CREDENTIALS = credentials('dockerhub_credentials')
        REPORT_EMAIL = "ahmedelmelegy3570@gmail.com"
    }

    stages {
        stage('Verify Branch') {
            steps {
                echo "$GIT_BRANCH"
            }
        }
        stage('Login to DockerHub') {
            steps {
                sh 'echo $DOCKERHUB_CREDENTIALS_PSW | docker login -u $DOCKERHUB_CREDENTIALS_USR --password-stdin'
            }
        }
        stage('Build image') {
            steps {
                catchError(buildResult: 'FAILURE', stageResult: 'FAILURE') {
                    sh(script: """
                        docker images
                        docker build -t ahmedelmelegy3570/app-multistage .
                    """)
                }
            }
        }
        stage('Push image') {
            steps {
                catchError(buildResult: 'FAILURE', stageResult: 'FAILURE') {
                    sh(script: """
                        docker push ahmedelmelegy3570/app-multistage
                    """)
                }
            }
        }
    }

    post {
        failure {
            mail to: REPORT_EMAIL,
                 subject: "Build Failed: ${env.JOB_NAME}",
                 body: "The build of ${env.JOB_NAME} has failed. Please check the Jenkins logs for more details."
        }
    }
}