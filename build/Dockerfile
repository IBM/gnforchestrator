# (C) Copyright IBM Corporation 2021.
# LICENSE: GPL-3.0 https://opensource.org/licenses/GPL-3.0

FROM quay.io/operator-framework/ansible-operator:v0.16.0

COPY requirements.yaml ${HOME}/requirements.yaml
RUN ansible-galaxy collection install -r ${HOME}/requirements.yaml \
 && chmod -R ug+rwx ${HOME}/.ansible

COPY watches.yaml ${HOME}/watches.yaml

COPY roles/ ${HOME}/roles/
COPY playbook.yaml ${HOME}/playbook.yaml
COPY ansible.cfg ${HOME}/ansible.cfg
ENV ANSIBLE_CONFIG=${HOME}/ansible.cfg
