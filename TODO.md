### **Checklist de Implementação - Value Objects**

- [x] Weight
- [x] Length
- [x] Longitude
- [x] Latitude
- [x] IPAddress
- [x] Slug
- [x] PortNumber
- [x] FileExtension
- [x] MIMEType
- [ ] Color
- [ ] Discount
- [ ] DayOfWeek
- [ ] BusinessHours
- [ ] Timezone

### **Conceito dos novos object**

A sua seleção de **Value Objects** reflete a maturidade na modelagem do domínio. Ao invés de usar tipos primitivos, estes objetos encapsulam regras de negócio, validações e comportamentos específicos, garantindo a integridade dos dados e tornando o código mais expressivo, seguro e fácil de manter.

---

#### **1. Medidas e Localização**

* **Weight (Peso):** Um Value Object que representa uma medida de massa com sua respectiva unidade. Ele garante que o valor seja sempre positivo e permite a conversão segura entre unidades como quilogramas, gramas ou libras, evitando erros em cálculos de logística ou estoque.
* **Length (Comprimento):** Similar ao Weight, este Value Object modela uma medida linear com sua unidade (metros, centímetros, pés). É fundamental para sistemas de estoque, logística ou dimensionamento de produtos, assegurando que os valores sejam válidos e consistentes.
* **Longitude e Latitude:** Estes Value Objects representam as coordenadas geográficas. Eles centralizam a validação para garantir que os valores numéricos estejam sempre dentro dos intervalos globais corretos (Latitude de -90 a 90; Longitude de -180 a 180), prevenindo dados incorretos em mapas ou cálculos de distância.

#### **2. Identificadores e Formatos**

* **Slug:** Um Value Object que representa uma string otimizada para ser um identificador amigável e único em URLs. Ele encapsula a lógica de formatação, como a conversão para letras minúsculas e a substituição de caracteres especiais por hífens, melhorando a usabilidade e o SEO.
* **IPAddress (Endereço IP):** Este objeto modela um endereço de rede. Ele centraliza a validação do formato do endereço (IPv4 ou IPv6), garantindo que os dados de log e rastreamento sejam sempre legítimos e confiáveis.
* **PortNumber (Número da Porta):** Essencial em ambientes de rede e arquiteturas de microserviços. Este Value Object garante que o valor da porta esteja no intervalo válido de 1 a 65535, evitando falhas de comunicação e protegendo contra portas inválidas.
* **FileExtension (Extensão de Arquivo):** Complementar ao MIMEType, este Value Object representa a extensão de um arquivo (.pdf, .xml). Ele garante a validade e a padronização do formato da extensão, facilitando o controle de tipos de arquivos em rotinas de upload.
* **MIMEType (Tipo de Mídia):** Um Value Object que define o tipo de conteúdo de um arquivo (por exemplo, image/jpeg). Ele é crucial para a segurança, pois valida o tipo de arquivo no momento do upload, protegendo o sistema de formatos maliciosos ou não suportados.
* **Color (Cor):** Um Value Object para representar cores (por exemplo, em formato Hexadecimal ou RGB). Ele encapsula a validação para garantir que a cor esteja em um formato correto e pode oferecer métodos para conversão entre formatos, sendo útil em interfaces visuais e dashboards.

#### **3. Lógica de Negócio e Tempo**

* **Discount (Desconto):** Este Value Object encapsula a lógica de um desconto, seja um valor fixo ou um percentual. Ele garante que os descontos sejam sempre valores válidos (não negativos) e centraliza a lógica de cálculo, prevenindo a duplicação de código e garantindo a consistência em operações financeiras.
* **DayOfWeek (Dia da Semana):** Representa um dia da semana de forma tipada e segura. Este objeto permite a criação de regras de negócio claras baseadas em dias específicos, como horários de operação ou agendamento de tarefas recorrentes.
* **BusinessHours (Horário Comercial):** Um Value Object mais complexo que define o período de operação de um negócio. Ele pode ser composto por outros objetos de tempo e encapsula a lógica para verificar se um determinado horário está dentro do expediente, o que é vital para sistemas de agendamento e suporte.
* **Timezone (Fuso Horário):** Este Value Object representa um fuso horário específico (e.g., "America/Sao_Paulo"). Ele garante que o nome do fuso seja válido e centraliza a lógica de conversão de tempo, simplificando o tratamento de datas em um ambiente global.
