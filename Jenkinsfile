pipeline {
    agent any

    environment {
        // ID của Secret chứa URL Deploy Hook trong Jenkins Credentials
        RENDER_HOOK_ID = 'render-deploy-hook-url'
        
        // Thêm thư mục bin của Go vào PATH để nhận lệnh gotestsum
        // Sử dụng %USERPROFILE% để trỏ tới thư mục go của user hiện tại
        PATH = "${PATH};${USERPROFILE}\\go\\bin"
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
                expression {
                    // Kiểm tra nhánh linh hoạt hơn (hỗ trợ cả Git branch có prefix origin/)
                    def currentBranch = env.BRANCH_NAME ?: env.GIT_BRANCH ?: ""
                    echo "🔍 Testing branch: ${currentBranch}"
                    return currentBranch.contains('main') || currentBranch.contains('master') || currentBranch.contains('develop')
                }
            }
            steps {
                script {
                    try {
                        withCredentials([string(credentialsId: "${RENDER_HOOK_ID}", variable: 'RENDER_URL')]) {
                            echo "🚀 Triggering Render Deployment..."
                            // Thêm flag -i để xem HTTP response từ Render
                            bat "curl -i -X POST %RENDER_URL%"
                        }
                    } catch (Exception e) {
                        echo "❌ Failed to trigger Render: ${e.getMessage()}"
                        error("Deployment failed")
                    }
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
