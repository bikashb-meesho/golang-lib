pipeline {
    agent any
    
    environment {
        GO_VERSION = '1.23'
        GOPATH = "${WORKSPACE}/go"
        PATH = "${GOPATH}/bin:/usr/local/go/bin:${env.PATH}"
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
                script {
                    env.GIT_COMMIT_SHORT = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
                    env.BRANCH_NAME = env.GIT_BRANCH.replaceAll('origin/', '')
                }
                echo "Building branch: ${env.BRANCH_NAME}"
                echo "Commit: ${env.GIT_COMMIT_SHORT}"
            }
        }
        
        stage('Setup Go') {
            steps {
                sh '''
                    go version
                    go env
                '''
            }
        }
        
        stage('Dependencies') {
            steps {
                echo 'Downloading dependencies...'
                sh '''
                    go mod download
                    go mod verify
                '''
            }
        }
        
        stage('Lint & Format Check') {
            steps {
                echo 'Running linters...'
                sh '''
                    # Check formatting
                    if [ -n "$(gofmt -l .)" ]; then
                        echo "The following files are not formatted:"
                        gofmt -l .
                        exit 1
                    fi
                    
                    # Run go vet
                    go vet ./...
                '''
            }
        }
        
        stage('Unit Tests') {
            steps {
                echo 'Running unit tests...'
                sh '''
                    go test -v -race -coverprofile=coverage.out ./...
                    go tool cover -func=coverage.out
                '''
            }
            post {
                always {
                    // Archive test results
                    sh 'go test -v -race ./... 2>&1 | tee test-results.txt || true'
                    archiveArtifacts artifacts: 'coverage.out,test-results.txt', allowEmptyArchive: true
                }
            }
        }
        
        stage('Coverage Check') {
            steps {
                echo 'Checking code coverage...'
                sh '''
                    # Extract coverage percentage
                    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
                    echo "Total coverage: ${COVERAGE}%"
                    
                    # Fail if coverage is below threshold (e.g., 70%)
                    THRESHOLD=70
                    if [ $(echo "$COVERAGE < $THRESHOLD" | bc) -eq 1 ]; then
                        echo "Coverage ${COVERAGE}% is below threshold ${THRESHOLD}%"
                        exit 1
                    fi
                '''
            }
        }
        
        stage('Build Verification') {
            steps {
                echo 'Verifying build...'
                sh '''
                    # Verify all packages can be built
                    go build ./...
                '''
            }
        }
        
        stage('Tag & Release') {
            when {
                branch 'main'
            }
            steps {
                script {
                    echo 'Main branch build successful'
                    echo 'Ready for tagging/release'
                    // Note: Actual tagging should be done manually or via separate release job
                    // to maintain semantic versioning control
                }
            }
        }
    }
    
    post {
        success {
            echo "✅ Build successful for ${env.BRANCH_NAME}"
            script {
                if (env.BRANCH_NAME == 'main') {
                    echo "Main branch validated - ready for version tagging"
                }
            }
        }
        failure {
            echo "❌ Build failed for ${env.BRANCH_NAME}"
        }
        always {
            cleanWs()
        }
    }
}

