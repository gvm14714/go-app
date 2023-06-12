pipeline {
   agent any
   environment {
    DOCKERHUB_CREDENTIALS=credentials('dockerhub_credentials')
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
        stage('build image') {
            steps {
                sh(script: """
                    docker images
                    docker build -t ahmedelmelegy3570/app-multistage .
                """)
                }
        }
        
        stage('push image') {
            steps {
                sh(script: """
                    docker push ahmedelmelegy3570/app-multistage
                """)
                }
        }
    }
    post {
        success {
            echo "Pipeline completed successfully"
            mail to: $REPORT_EMAIL,
            subject: "go-app-intern",
            body: "Pipeline completed successfully"
        }
        failure {
            echo "Pipeline failed"
            mail to: $REPORT_EMAIL,
            subject: "go-app-intern",
            body: "Pipeline failed ... Please Check logs"
        }

    }
}