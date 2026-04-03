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
                sh 'go mod download'
                sh 'go install gotest.tools/gotestsum@latest'
            }
        }

        stage('Unit Tests') {
            steps {
                // Chạy Unit Test và tạo báo cáo JSON
                sh 'go test -json -v ./... > unit_report.json'
                // Hiển thị kết quả tóm tắt
                sh 'gotestsum --format short-verbose < unit_report.json'
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
                    // Bật server background
                    sh 'go run main.go &'
                    // Đợi server sẵn sàng
                    sh "timeout 30s sh -c 'until curl -s http://127.0.0.1:8080/ping; do sleep 1; done'"
                    // Chạy Automation Test
                    sh 'go test -json -v -tags automation ./tests/automation/... > auto_report.json || true'
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
                // Sử dụng withCredentials để lấy URL bảo mật từ Jenkins
                withCredentials([string(credentialsId: "${RENDER_HOOK_ID}", variable: 'RENDER_URL')]) {
                    echo "🚀 Triggering Render Deployment..."
                    sh "curl -X POST \"${RENDER_URL}\""
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
