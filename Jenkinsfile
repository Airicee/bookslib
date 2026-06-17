pipeline {
    agent any

    environment {
        REGISTRY          = "bookslib-local"
        AUTH_SERVICE      = "bookslib-auth-service"
        BOOKS_SERVICE     = "bookslib-books-service"
        REVIEWS_SERVICE   = "bookslib-reviews-service"
        FRONTEND_SERVICE  = "bookslib-frontend"
    }

    stages {
        stage('1. Environment Check') {
            steps {
                echo '=== STAGE: CHECKING TOOLS VERSION ==='
                sh 'docker version'
                sh 'docker compose version'
            }
        }

        stage('2. Static Application Security Testing (SAST)') {
            steps {
                echo '=== STAGE: RUNNING SECURITY SOURCE CODE SCANNING ==='
                echo 'Executing Permanently Installed Trivy...'
                sh 'trivy fs --severity HIGH,CRITICAL --exit-code 0 .'
            }
        }

        stage('3. Build Microservices Images') {
            steps {
                echo '=== STAGE: BUILDING DOCKER IMAGES WITH COMPOSE ==='
                sh 'docker compose build --no-cache'
            }
        }

        stage('4. Verify Images') {
            steps {
                echo '=== STAGE: VERIFYING IMAGES GENERATION ==='
                sh 'docker images | grep bookslib'
            }
        }

        stage('5. Deploy / Run Application') {
            steps {
                echo '=== STAGE: DEPLOYING MICROSERVICES ==='
                sh 'docker compose down && docker compose up -d'
            }
        }
    }

    post {
        success {
            echo 'Pipeline build successfully completed!'
        }
        failure {
            echo 'Pipeline failed. Please check the console logs for debugging.'
        }
    }
}