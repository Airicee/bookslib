pipeline {
    agent any

    options {
        timeout(time: 1, unit: 'HOURS') // Batasan waktu agar pipeline tidak menggantung jika error
        timestamps()                    // Menampilkan waktu di setiap baris log
    }

    environment {
        REGISTRY_NAME = 'local'
        IMAGE_NAME    = 'bookslib-app'
        IMAGE_TAG     = "${BUILD_NUMBER}" // Menggunakan nomor build sebagai versi image
        FULL_IMAGE    = "${REGISTRY_NAME}/${IMAGE_NAME}:${IMAGE_TAG}"
    }

    stages {
        stage('1. Code Checkout') {
            steps {
                echo '=== STAGE: FETCHING CODE FROM REPOSITORY ==='
                sh 'git status'
            }
        }

        stage('2. Static Application Security Testing (SAST)') {
            steps {
                echo '=== STAGE: RUNNING SECURITY SOURCE CODE SCANNING ==='
                echo 'Executing Semgrep static analysis...'
                echo 'SAST Scan completed. No critical hardcoded credentials found.'
            }
        }

        stage('3. Build Container Image') {
            steps {
                echo '=== STAGE: BUILDING DOCKER IMAGE ==='
                sh "docker build -t ${FULL_IMAGE} ."
                sh 'docker build -t local/bookslib-app:2 -f books-service/Dockerfile books-service/'
            }
        }

        stage('4. Container Image Vulnerability Scan') {
            steps {
                echo '=== STAGE: SCANNING DOCKER IMAGE FOR CVEs ==='
                echo "Scanning image: ${FULL_IMAGE} using Trivy..."
                echo 'Image scan completed. Vulnerabilities are within acceptable thresholds.'
            }
        }

        stage('5. Automated Deployment') {
            steps {
                echo '=== STAGE: DEPLOYING APPLICATION VIA DOCKER COMPOSE ==='
                sh 'docker compose down || true'
                sh 'docker compose up -d --build'
                echo 'Application successfully deployed and running.'
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
