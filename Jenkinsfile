pipeline {
    agent {
        node {
            label "psi_rhel7_openshift311"
        }
    }

    libraries {
        lib('fh-pipeline-library')
        lib('qe-pipeline-library')
    }
    
    environment {
        GOPATH = "${env.WORKSPACE}/"
        PATH = "${env.PATH}:${env.WORKSPACE}/bin:/usr/local/go/bin"
        GOOS = "linux"
        GOARCH = "amd64"
        CGO_ENABLED = 0
        OPERATOR_NAME = "unifiedpush-operator"
        OPENSHIFT_PROJECT_NAME = "unifiedpush"
    }

    options {
        checkoutToSubdirectory("src/github.com/aerogear/unifiedpush-operator")
    }

    stages {

        stage("Trust"){
            steps{
                enforceTrustedApproval('aerogear')
            }
            post{
                failure{
                    echo "====++++'Trust' execution failed++++===="
                    echo "You are not authorized to run this job"
                }
            }
        }

        stage("Run e2e test") {
            steps {
                script {
                    testKubernetesOperator(
                        clonedRepositoryPath: "src/github.com/aerogear/unifiedpush-operator",
                        openshiftProjectName: "${env.OPENSHIFT_PROJECT_NAME}",
                        operatorContainerImageCandidateName: "quay.io/aerogear/${env.OPERATOR_NAME}:candidate-${env.BRANCH_NAME}",
                        operatorContainerImageName: "quay.io/aerogear/${env.OPERATOR_NAME}:${env.BRANCH_NAME}",
                        operatorContainerImageNameLatest: "quay.io/aerogear/${env.OPERATOR_NAME}:latest",
                        containerRegistryCredentialsId: "quay-aerogear-bot"
                    )
                }
            }
        }
    }
    post {
        always{
            script {
                sh """
                oc delete project ${env.OPENSHIFT_PROJECT_NAME}
                rm -rf ${env.CLONED_REPOSITORY_PATH}
                """
            }
        }
        failure {
            mail(
                to: 'psturc@redhat.com',
                subject: 'UnifiedPush Operator build failed',
                body: "See the pipeline here: ${env.RUN_DISPLAY_URL}"
            )
        }
    }
}