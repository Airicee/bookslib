pipeline {
    agent any

    options {
        timeout(time: 1, unit: 'HOURS')
        timestamps()
    }

    environment {
        TARGET_IMAGE = 'local/bookslib-auth:latest'
    }

    stages {
        stage('1. Code Checkout') {
            steps {
                echo '=== STAGE: FETCHING CODE FROM REPOSITORY ==='
                checkout scm
                sh 'git status'
            }
        }

        stage('2. Static Application Security Testing (SAST)') {
             steps {
                echo '=== STAGE: RUNNING SECURITY SOURCE CODE SCANNING ==='
                echo 'Downloading and Executing Trivy via Script...'
                sh '''
                    curl -Lo trivy_0.49.1_Linux-64bit.tar.gz https://github.com/aquasecurity/trivy/releases/download/v0.49.1/trivy_0.49.1_Linux-64bit.tar.gz
                    tar -xzf trivy_0.49.1_Linux-64bit.tar.gz trivy
                    ./trivy fs --severity HIGH,CRITICAL --exit-code 0 .
                    rm trivy trivy_0.49.1_Linux-64bit.tar.gz
                '''
            }
        }

        stage('3. Build Container Image') {
            steps {
                echo '=== STAGE: BUILDING DOCKER IMAGES FOR ALL SERVICES ==='
                
                echo 'Building Auth Service...'
                sh 'docker build -t local/bookslib-auth:latest ./auth-service'
                
                echo 'Building Books Service...'
                sh 'docker build -t local/bookslib-books:latest ./books-service'
                
                echo 'Building Reviews Service...'
                sh 'docker build -t local/bookslib-reviews:latest ./reviews-service'
                
                echo 'Building Frontend...'
                sh 'docker build -t local/bookslib-frontend:latest ./frontend'
            }
        }

        stage('4. Container Image Vulnerability Scan') {
            steps {
                echo '=== STAGE: SCANNING DOCKER IMAGE FOR CVEs ==='
                echo "Scanning image: ${TARGET_IMAGE} using Trivy..."
                sh "trivy image --severity CRITICAL --exit-code 0 ${TARGET_IMAGE}"
            }
        }

        stage('5. Automated Deployment') {
            steps {
                echo '=== STAGE: DEPLOYING APPLICATION VIA DOCKER COMPOSE ==='
                sh 'docker compose down || true'
                sh 'docker compose up -d --build'
            }
        }
    }

    post {
        always {
            echo '=== POST BUILD: CLEANING UP WORKSPACE ==='
            sh 'docker image prune -f || true'
        }
        success {
            echo "✅ PIPELINE SUCCESSFUL: Build #${BUILD_NUMBER} telah berhasil dideploy dengan aman."
        }
        failure {
            echo "❌ PIPELINE FAILED: Terjadi kesalahan pada Build #${BUILD_NUMBER}. Silakan periksa log di atas."
        }
    }
}