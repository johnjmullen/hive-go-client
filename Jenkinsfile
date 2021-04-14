pipeline {
  agent {
    docker {
      image 'hive/package-builder-ci:8.3'
      args '-u root:root --network ci_ci'
    }
  }
  stages {
    stage('Package') {
      steps {
          sh 'add-apt-repository ppa:hnakamur/golang-1.13'
          sh '/package.sh'
      }
    }
  }
  post {
     success {
            archiveArtifacts artifacts: 'package-urls.txt', fingerprint: true
        }
    }
}
