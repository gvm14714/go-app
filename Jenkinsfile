pipeline {
    agent any
    environment {
        // Declaring environment variable for report email
        REPORT_EMAIL = "gvm14714@gmail.com"
    }

    stages {
        stage('Verify Branch') {
            steps {
                script {
                    // Checking for the BRANCH_NAME (Multibranch Pipeline) or GIT_BRANCH
                    def branchName = env.BRANCH_NAME ?: env.GIT_BRANCH ?: 'Unknown'
                    echo "Branch: ${branchName}"
                }
            }
        }
        stage('Login to DockerHub') {
            steps {
                // Using withCredentials to handle DockerHub credentials
                // withCredentials([usernamePassword(credentialsId: 'dockerhub_credentials', usernameVariable: 'DOCKERHUB_USER', passwordVariable: 'DOCKERHUB_PSW')]) {
                //     sh 'echo $DOCKERHUB_PSW | docker login -u $DOCKERHUB_USER --password-stdin'
                // }
            }
        }
        stage('Build image') {
            steps {
                catchError(buildResult: 'FAILURE', stageResult: 'FAILURE') {
                    sh(script: """
                        docker images
                        docker build -t gym14714/us .
                    """)
                }
            }
        }
        stage('Push image') {
            steps {
                catchError(buildResult: 'FAILURE', stageResult: 'FAILURE') {
                    sh(script: """
                        docker push gym14714/us
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
