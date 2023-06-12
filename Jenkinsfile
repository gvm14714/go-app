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
            mail bcc: '', body: "<b>Example</b><br>Project: ${env.JOB_NAME} <br>Build Number: ${env.BUILD_NUMBER} <br> URL de build: ${env.BUILD_URL}", cc: '', charset: 'UTF-8', from: '', mimeType: 'text/html', replyTo: '', subject: "ERROR CI: Project name -> ${env.JOB_NAME}", to: "ahmedelmelegy3570@gmail.com";
        }
        failure {
            echo "Pipeline failed"
            mail bcc: '', body: "<b>Example</b><br>Project: ${env.JOB_NAME} <br>Build Number: ${env.BUILD_NUMBER} <br> URL de build: ${env.BUILD_URL}", cc: '', charset: 'UTF-8', from: '', mimeType: 'text/html', replyTo: '', subject: "ERROR CI: Project name -> ${env.JOB_NAME}", to: "ahmedelmelegy3570@gmail.com";
        }

    }
}