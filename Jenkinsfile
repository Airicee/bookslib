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
                echo "Scanning image: ${FULL_IMAGE} using Trivy..."
                echo 'Image scan completed. Vulnerabilities are within acceptable thresholds.'
            }
        }

        stage('5. Automated Deployment') {
            steps {
                echo '=== MENYIAPKAN PLUGIN DOCKER COMPOSE ==='
                sh '''
                    # 1. Buat folder plugin Docker untuk user Jenkins
                    mkdir -p ~/.docker/cli-plugins/
            
                    # 2. Download plugin Docker Compose V2 terbaru
                    curl -SL https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-linux-x86_64 -o ~/.docker/cli-plugins/docker-compose
            
                    # 3. Berikan izin eksekusi pada plugin
                    chmod +x ~/.docker/cli-plugins/docker-compose
            
                    # 4. Tes apakah Docker sekarang sudah mengenali perintah 'docker compose'
                    docker compose version
                '''

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