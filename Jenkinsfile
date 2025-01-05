pipeline {
    agent none
    stages {
        stage('Setup Environment') {
            steps {
                script {
                    if (env.BRANCH_NAME == 'master') {
                        agent_label = 'master-agent'
                        docker_compose_file = 'docker-compose.prod.yml'
                        credentialsId = 'backend-prod'
                    } else if (env.BRANCH_NAME == 'develop') {
                        agent_label = 'develop-agent'
                        docker_compose_file = 'docker-compose.dev.yml'
                        credentialsId = 'backend-dev'
                    } else {
                        error "Branch ${env.BRANCH_NAME} is not configured!"
                    }
                }
            }
        }
        
        stage('Build and Deploy') {
            agent { label agent_label }
            steps {
                withCredentials([file(credentialsId: credentialsId, variable: 'env_file')]) {
                    // Set environment variables
                    sh "cat ${env_file} > .env"
                    // Start services
                    sh "docker compose -f ${docker_compose_file} up -d --build"
                }
            }
        }
    }

    post {
        always {
            echo "Pipeline finished for branch: ${env.BRANCH_NAME}"
        }
    }
}