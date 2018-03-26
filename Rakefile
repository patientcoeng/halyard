namespace :image do
  task :prep do
    sh "rm -f halyard"
    sh "govendor sync"
    sh "go build ."
  end
  task :build => :prep do
    sh "sudo docker build -t registry.patientco.engineering/halyard ."
  end
  task :push_prep do
    sh "aws ecr get-login --no-include-email --region us-east-1 --profile=docker | sed 's#https://.*#https://registry.patientco.engineering#' | bash"
    sh "aws ecr create-repository --profile=docker --repository-name halyard --region=us-east-1 || true"
  end
  task :push => [:push_prep, :build] do
    sh "sudo docker tag registry.patientco.engineering/halyard registry.patientco.engineering/halyard:$(git rev-parse --short HEAD) && \
        sudo docker push registry.patientco.engineering/halyard:$(git rev-parse --short HEAD)"
  end
  task :export => :build do
    sh "sudo docker save -o halyard.tar registry.patientco.engineering/halyard"
  end
  task :test do
    sh "go test -v ./..."
  end
end
