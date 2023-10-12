[
  {
    kind: 'pipeline',
    type: 'docker',
    name: 'default',
    clone: {
      disable: true,
    },
    workspace: {
      path: '/workspaces/hive-go-client',
    },
    steps: [
      {
        name: 'clone',
        image: 'alpine/git',
        commands: [
          'cp -r /home/admin1/.ssh /root/',
          'touch /root/.ssh/known_hosts',
          'chmod 600 /root/.ssh/known_hosts',
          'chown -R root:root /root/.ssh',
          'ssh-keyscan -H github.com > ~/.ssh/known_hosts',
          'git config --add --global url."git@github.com:".insteadOf https://github.com',
          'git clone git@github.com:hive-io/hive-go-client.git .',
          'git checkout $DRONE_COMMIT',
        ],
        volumes: [
          {
            name: 'sshkey',
            path: '/home/admin1/.ssh',
          },
          {
            name: 'root-dot-ssh',
            path: '/root/.ssh',
          },
        ],
      },
      {
        name: 'build image',
        image: 'plugins/docker',
        settings: {
          debug: true,
        },
        privileged: true,
        volumes: [
          {
            name: 'dockersock',
            path: '/var/run/docker.sock',
          },
        ],
        commands: [
          'docker build -f .devcontainer/Dockerfile -t hive-go-client-${DRONE_COMMIT_BRANCH,,} .',
          'docker run --name hive-go-client-${DRONE_COMMIT_BRANCH,,}-${DRONE_COMMIT_SHA:0:8} hive-go-client-${DRONE_COMMIT_BRANCH,,}',
        ],
        dry_run: true,
      },
      {
        name: 'run test',
        pull: 'never',
        image: 'hive-go-client-${DRONE_COMMIT_BRANCH,,}',
        privileged: true,
        volumes: [
          {
            name: 'sshkey',
            path: '/home/admin1/.ssh',
          },
          {
            name: 'root-dot-ssh',
            path: '/root/.ssh',
          },
        ],
        commands: [
          '/usr/bin/testenv.sh',
        ],
        when: {
          status: [
            'success',
          ],
        },
      },
      {
        name: 'build package',
        pull: 'never',
        image: 'hive-go-client-${DRONE_COMMIT_BRANCH,,}',
        privileged: true,
        volumes: [
          {
            name: 'sshkey',
            path: '/home/admin1/.ssh',
          },
          {
            name: 'root-dot-ssh',
            path: '/root/.ssh',
          },
        ],
        commands: [
          'WORK_DIR=$(mktemp -d)',
          'cp -ar . $WORK_DIR',
          'cd $WORK_DIR',
          'debian/update-changelog.sh',
          'dpkg-buildpackage -b',
          '/usr/bin/publish-aptly.sh $WORK_DIR $DRONE_COMMIT_BRANCH',
        ],
        when: {
          status: [
            'success',
          ],
        },
      },
      {
        name: 'cleanup',
        image: 'plugins/docker',
        volumes: [
          {
            name: 'dockersock',
            path: '/var/run/docker.sock',
          },
        ],
        privileged: true,
        failure: 'ignore',
        commands: [
          'docker container stop hive-go-client-${DRONE_COMMIT_BRANCH,,}-${DRONE_COMMIT_SHA:0:8}',
          'docker container rm -f hive-go-client-${DRONE_COMMIT_BRANCH,,}-${DRONE_COMMIT_SHA:0:8}',
        ],
        when: {
          status: [
            'failure',
            'success',
          ],
        },
      },
    ],
    trigger: {
      event: null,
      include: [
        'push',
        'tag',
      ],
      exclude: [
        'pull_request',
      ],
    },
    volumes: [
      {
        name: 'dockersock',
        host: {
          path: '/var/run/docker.sock',
        },
      },
      {
        name: 'sshkey',
        host: {
          path: '/home/admin1/.ssh',
        },
      },
      {
        name: 'root-dot-ssh',
        temp: {},
      },
    ],
  },
]