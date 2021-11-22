pipeline {
  agent any
  stages {
    stage('parallel_tests') {
      parallel {
        stage('test1') {
          environment {
            PROD = '1'
          }
          steps {
            sh 'docker image ls'
          }
        }

        stage('test2') {
          environment {
            DEV = '0'
          }
          steps {
            echo 'test2'
          }
        }

      }
    }

    stage('test_child') {
      steps {
        timeout(time: 10) {
          sh 'ls'
        }

      }
    }

  }
  environment {
    DEV_MODE = '1'
  }
}