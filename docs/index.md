# Локальное развертывание
1. Клонируйте репозиторий:
```bash
git clone https://github.com/fle4a/transaction-system.git your-repo
```

2. Перейдите в директорию проекта:
```bash
cd your-repo
```

3. Переопределить настройки в соответствии с рабочим окружением
```bash
find . -type f -path "*/src/configs/local.yaml.tmpl" -exec sh -c 'cp "$0" "$(dirname "$0")/local.yaml"' {} \;
```

4. Запустить приложение
```bash
docker compose up
```

5. Теперь сервис доступен по адресу `localhost:8080`

6. Предварительный запуск для создания кошельков, сохраняемых в файл uuids.txt
```bash
cd tests && go run && cd ..
```
