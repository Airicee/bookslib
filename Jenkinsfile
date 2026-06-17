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
                // Tahap ini berjalan aman dengan kondisi socket masih terkunci (660)
                sh 'trivy fs --severity HIGH,CRITICAL --exit-code 0 .'
            }
        }

        stage('3. Build Microservices Images') {
            steps {
                echo '=== JIT SECURITY: TEMPORARILY OPENING DOCKER SOCKET FOR BUILD ==='
                // Membuka pintu socket secara berkala tepat saat dibutuhkan
                sh 'chmod 777 /var/run/docker.sock || true'

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
                // Menjalankan container web dan microservices ke background
                sh 'docker compose down && docker compose up -d'
            }
        }
    }

    post {
        always {
            echo '=== CLEANUP SECURITY: RESTRICTING DOCKER SOCKET PERMISSION IMMEDIATELY ==='
            // Menutup dan mengunci kembali pintu socket ke mode aman (660) setelah pipeline selesai
            sh 'chmod 660 /var/run/docker.sock || true'
        }
        success {
            echo '🎉 Pipeline build and deployment successfully completed!'
        }
        failure {
            echo '❌ Pipeline failed. Docker socket safely locked down.'
        }
    }
}