pipeline {
   agent any
   environment {
    DOCKERHUB_CREDENTIALS=credentials('dockerhub_credentials')
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
        stage('build and push image') {
            steps {
                sh(script: """
                    docker images
                    docker build -t ahmedelmelegy3570/app-multistage .
                    docker push ahmedelmelegy3570/app-multistage
                """)
                }
        }
   }
}