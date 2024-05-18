module.exports = function (plop) {
  plop.setHelper('kubecommand', function (cluster) {
    return {
      'mmz-nodes': 'kubectl',
      'macberengar (deprecated)': 'kubectl --context macberengar',
      none: 'echo'
    }[cluster]
  })

  plop.setGenerator('basics', {
    description: 'Creates a skeleton for a webapplication with Sveltekit frontend and FastAPI python backend',
    prompts: [
      {
        name: 'name',
        type: 'input',
        message: 'give the name of your project'
      },
      {
        name: 'cluster',
        type: 'list',
        choices: ['mmz-nodes', 'macberengar (deprecated)', 'none'],
        message: 'select your k8s cluster'
      },
      {
        name: 'namespace',
        type: 'input',
        message: 'give the namespace of your project',
        default: 'default'
      },
      {
        name: 'gitprovider',
        type: 'list',
        choices: ['manual', 'gitlab', 'github'],
        message: 'select your git repository provider'
      }
    ],
    actions: [
      {
        type: 'add',
        path: 'deploy.sh',
        templateFile: '.deploy.sh.hbs'
      },
      {
        type: 'add',
        path: 'Deployment.yaml',
        templateFile: '.Deployment.yaml.hbs'
      },
      {
        type: 'add',
        path: 'Ingress.yaml',
        templateFile: '.Ingress.yaml.hbs'
      },
      {
        type: 'add',
        path: 'Services.yaml',
        templateFile: '.Services.yaml.hbs'
      }
    ]
  })
}
