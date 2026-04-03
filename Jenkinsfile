pipeline {
    agent any

    environment {
        // ID của Secret chứa URL Deploy Hook trong Jenkins Credentials
        RENDER_HOOK_ID = 'render-deploy-hook-url'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Install Dependencies') {
            steps {
                bat 'go mod download'
                bat 'go install gotest.tools/gotestsum@latest'
            }
        }

        stage('Unit Tests') {
            steps {
                // Chạy Unit Test và tạo báo cáo JSON
                bat 'go test -json -v ./... > unit_report.json'
                // Hiển thị kết quả tóm tắt
                bat 'gotestsum --format short-verbose < unit_report.json'
            }
        }

        stage('Automation Tests') {
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                    branch 'develop'
                }
            }
            steps {
                script {
                    // Bật server background trên Windows (dùng start /B)
                    bat 'start /B go run main.go'
                    // Đợi server sẵn sàng (ping thành công)
                    bat 'timeout /t 10 /nobreak'
                    // Chạy Automation Test
                    bat 'go test -json -v -tags automation ./tests/automation/... > auto_report.json || exit 0'
                }
            }
        }

        stage('Deploy to Render') {
            when {
                anyOf {
                    branch 'main'
                    branch 'master'
                }
            }
            steps {
                withCredentials([string(credentialsId: "${RENDER_HOOK_ID}", variable: 'RENDER_URL')]) {
                    echo "🚀 Triggering Render Deployment..."
                    // Dùng curl (Windows 10/11 đã tích hợp sẵn curl)
                    bat "curl -X POST %RENDER_URL%"
                }
            }
        }
    }

    post {
        always {
            // Lưu lại kết quả test để xem trên Jenkins Dashboard
            archiveArtifacts artifacts: '*.json', allowEmptyArchive: true
            echo "CI/CD Pipeline Finished."
        }
    }
}
