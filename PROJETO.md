## Documento Técnico Detalhado: PlanSmith

### 1. Introdução

O PlanSmith é uma aplicação de linha de comando (CLI) desenvolvida em Go, que atua como uma ferramenta de assistência para o planejamento e gerenciamento de projetos de software. A aplicação se propõe a resolver um problema comum no desenvolvimento de software: a dificuldade de traduzir ideias de projetos, muitas vezes descritas de forma informal em documentos de texto, em um plano de projeto estruturado e acionável.

O objetivo principal do PlanSmith é automatizar a criação de quadros Kanban no Trello a partir de uma descrição de projeto em Markdown. Para isso, a aplicação utiliza inteligência artificial para analisar o documento de entrada, gerar uma visão do produto, épicos, estórias de usuário e tarefas técnicas, e, por fim, criar um quadro no Trello com todos esses elementos.

O projeto foi desenvolvido com foco na interatividade e na revisão humana. O usuário pode acompanhar todo o processo de criação do plano, revisar e aprovar cada etapa, garantindo que o resultado final esteja de acordo com as suas expectativas.

### 2. Arquitetura do Sistema

O PlanSmith foi projetado com uma arquitetura modular e orientada a eventos, utilizando a biblioteca `bubbletea` para a construção da sua interface de terminal (TUI). A arquitetura pode ser dividida nos seguintes componentes principais:

![Arquitetura do PlanSmith](https://i.imgur.com/9Y7b3yG.png)

**Fluxo de Dados:**

1.  O **Usuário** interage com a **TUI**, fornecendo comandos e arquivos de entrada.
2.  A **TUI** envia os dados de entrada para o **Agente Smith**.
3.  O **Agente Smith** utiliza o **Loader de Prompts** para carregar os modelos de prompt apropriados.
4.  O **Agente Smith** envia os prompts para o **Executor de IA**.
5.  O **Executor de IA** se comunica com a **API de IA** (Gemini, OpenAI, etc.) e retorna a resposta para o **Agente Smith**.
6.  O **Agente Smith** processa a resposta da IA e a envia de volta para a **TUI**.
7.  A **TUI** exibe o plano gerado para o **Usuário** para revisão.
8.  Após a aprovação do **Usuário**, a **TUI** envia o plano para o **Cliente Trello**.
9.  O **Cliente Trello** se comunica com a **API do Trello** para criar o quadro Kanban.
10. O **Gerenciador de Estado** é responsável por salvar e carregar o estado da aplicação, incluindo o plano do projeto e as configurações do Trello.

### 3. Componentes Detalhados

#### 3.1. `pkg/tui` (Terminal User Interface)

A TUI é o coração da interação com o usuário. Ela é responsável por:

*   **Apresentar a interface de chat:** A TUI utiliza a biblioteca `bubbletea` para criar uma interface de chat interativa, onde o usuário pode inserir comandos e receber feedback da aplicação.
*   **Gerenciar os estados da aplicação:** A TUI é uma máquina de estados, onde cada estado representa uma etapa do fluxo de trabalho (ex: menu principal, entrada de arquivo, revisão do plano, etc.).
*   **Lidar com a entrada do usuário:** A TUI captura os comandos do usuário, os processa e os envia para os componentes apropriados.
*   **Exibir o feedback da aplicação:** A TUI exibe as mensagens da aplicação, o plano gerado, os indicadores de carregamento e as mensagens de erro.

**Tecnologias Utilizadas:**

*   **`charmbracelet/bubbletea`:** Framework para a construção de TUIs em Go.
*   **`charmbracelet/lipgloss`:** Biblioteca para estilização de texto no terminal.
*   **`charmbracelet/bubbles`:** Coleção de componentes `bubbletea` (spinner, textinput, etc.).

#### 3.2. `pkg/smith` (Agente Smith)

O Agente Smith é o cérebro da aplicação. Ele é responsável por:

*   **Orquestrar o fluxo de trabalho:** O Agente Smith coordena a interação entre a TUI, o Executor de IA e o Cliente Trello.
*   **Gerar o plano do projeto:** O Agente Smith utiliza o Executor de IA para gerar a visão do produto, os épicos, as estórias de usuário e as tarefas técnicas.
*   **Processar as respostas da IA:** O Agente Smith processa as respostas da IA, as estrutura em um plano de projeto e as envia para a TUI.

#### 3.3. `pkg/ai` (Executor de IA)

O Executor de IA é uma abstração para interagir com diferentes provedores de IA. Ele é responsável por:

*   **Enviar os prompts para a IA:** O Executor de IA envia os prompts para a API de IA e aguarda a resposta.
*   **Lidar com diferentes provedores de IA:** O Executor de IA suporta diferentes provedores de IA (Gemini, OpenAI, Qwen), permitindo que o usuário escolha o provedor de sua preferência.

#### 3.4. `pkg/trello` (Cliente Trello)

O Cliente Trello é responsável por interagir com a API do Trello. Ele é responsável por:

*   **Criar o quadro Kanban:** O Cliente Trello cria um novo quadro no Trello com as listas padrão (Epics, User Stories, To Do, Doing, Done).
*   **Popular o quadro:** O Cliente Trello cria os cartões no quadro para cada épico, estória de usuário e tarefa.
*   **Vincular os cartões:** O Cliente Trello cria os vínculos entre os cartões para representar as dependências entre as tarefas.

#### 3.5. `pkg/state` (Gerenciador de Estado)

O Gerenciador de Estado é responsável por salvar e carregar o estado da aplicação. Ele é responsável por:

*   **Salvar o plano do projeto:** O Gerenciador de Estado salva o plano do projeto em um arquivo JSON.
*   **Salvar o estado da aplicação:** O Gerenciador de Estado salva o estado da aplicação, incluindo o ID e a URL do quadro do Trello.

### 4. Fluxos de Trabalho Detalhados

#### 4.1. Fluxo de Trabalho de Novo Projeto

![Fluxo de Trabalho de Novo Projeto](https://i.imgur.com/1c2g3hW.png)

#### 4.2. Fluxo de Trabalho de Adicionar Funcionalidade

![Fluxo de Trabalho de Adicionar Funcionalidade](https://i.imgur.com/4j5f6gH.png)

### 5. Conclusão

O PlanSmith é uma ferramenta poderosa que automatiza o processo de criação de planos de projeto, permitindo que as equipes de desenvolvimento de software se concentrem no que realmente importa: o desenvolvimento do produto. A aplicação é flexível, interativa e fácil de usar, e pode ser facilmente estendida para suportar novos provedores de IA e novas funcionalidades.

**Trabalhos Futuros:**

*   **Edição do plano na TUI:** Permitir que o usuário edite o plano gerado diretamente na TUI antes de criar o quadro no Trello.
*   **Suporte a outros gerenciadores de projetos:** Adicionar suporte a outros gerenciadores de projetos, como o Jira e o Asana.
*   **Melhorias na IA:** Aprimorar os prompts da IA para gerar planos de projeto ainda mais precisos e detalhados.
*   **Visualização de dependências:** Adicionar uma visualização gráfica das dependências entre as tarefas.
*   **Internacionalização:** Adicionar suporte a outros idiomas.

### 6. Lógica de Negócio Principal

A lógica de negócio principal do PlanSmith reside na sua capacidade de transformar uma descrição de projeto em Markdown em um plano de projeto estruturado e acionável. Este processo é orquestrado pelo **Agente Smith** e utiliza uma cadeia de prompts de IA para gerar os diferentes elementos do plano.

O processo pode ser dividido nas seguintes etapas:

1.  **Geração da Visão do Produto:**
    *   O processo começa com o prompt `01_generate_vision.md`. Este prompt instrui a IA a atuar como um Product Owner e a gerar um nome de projeto, uma visão de produto e uma lista de épicos de alto nível a partir do arquivo Markdown fornecido pelo usuário.
    *   A resposta da IA é um objeto JSON que é então processado pelo **Agente Smith** e exibido na TUI para revisão do usuário.

2.  **Geração de Estórias de Usuário:**
    *   Após a aprovação da visão do produto, o **Agente Smith** itera sobre cada épico e utiliza o prompt `02_generate_stories.md`.
    *   Este prompt instrui a IA a atuar como um Product Owner e a quebrar cada épico em estórias de usuário, seguindo o formato "Como um [USUÁRIO], eu quero [FUNCIONALIDADE], para que [VALOR]".
    *   A resposta da IA é um objeto JSON contendo uma lista de estórias de usuário, que são então adicionadas ao plano do projeto.

3.  **Geração de Tarefas Técnicas:**
    *   Com as estórias de usuário geradas, o **Agente Smith** itera sobre cada estória e utiliza o prompt `03_generate_tasks.md`.
    *   Este prompt instrui a IA a atuar como um Lead Developer e a quebrar cada estória de usuário em tarefas técnicas granulares.
    *   O prompt também instrui a IA a identificar as dependências entre as tarefas, o que é crucial para a criação de um plano de projeto realista.
    *   A resposta da IA é um objeto JSON contendo uma lista de tarefas, que são então adicionadas ao plano do projeto.

4.  **Criação do Quadro no Trello:**
    *   Após a geração de todo o plano, o usuário pode optar por criar um quadro no Trello.
    *   O **Cliente Trello** é então utilizado para criar o quadro, as listas, os cartões para cada épico, estória de usuário e tarefa, e os vínculos entre os cartões para representar as dependências.

**Cadeia de Prompts (Prompt Chaining):**

A utilização de uma cadeia de prompts é um dos principais diferenciais do PlanSmith. Em vez de enviar um único prompt para a IA e esperar que ela gere todo o plano de uma só vez, o PlanSmith divide o processo em etapas, utilizando um prompt específico para cada etapa.

Esta abordagem oferece várias vantagens:

*   **Melhor controle sobre o processo:** Permite que o usuário revise e aprove cada etapa do processo, garantindo que o resultado final esteja de acordo com as suas expectativas.
*   **Resultados mais precisos:** Permite que a IA se concentre em uma única tarefa de cada vez, o que resulta em resultados mais precisos e detalhados.
*   **Maior flexibilidade:** Permite que o usuário intervenha no processo e faça ajustes, se necessário.

**Gerenciamento de Estado:**

O **Gerenciador de Estado** desempenha um papel crucial na lógica de negócio do PlanSmith. Ele é responsável por salvar o plano do projeto em um arquivo JSON, o que permite que o usuário:

*   **Continue o trabalho de onde parou:** O usuário pode fechar a aplicação e continuar o trabalho de onde parou, sem perder o progresso.
*   **Adicione novas funcionalidades a projetos existentes:** O usuário pode carregar um plano de projeto existente e adicionar novas funcionalidades a ele, utilizando o fluxo de trabalho de "Adicionar Funcionalidade".

Em resumo, a lógica de negócio principal do PlanSmith é uma combinação de uma cadeia de prompts de IA bem definida, um gerenciamento de estado robusto e uma interação contínua com o usuário, o que resulta em uma ferramenta poderosa e flexível para o planejamento e gerenciamento de projetos de software.