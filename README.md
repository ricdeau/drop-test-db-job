# drop-test-db-job

**Описание**

Запускает джоб, удаляющий базы данных, название которых соответствует 
определенному шаблону, и которые старше указанного времени хранения.

* Ищет базы данных по regex шаблону *^\d{10}-* в названии.
* Первые 10 цифр в названии бд являются числовым форматом UTC.
* Конфигурируемые переменные среды:
    - *DB_TYPE* - тип базы данных
        - postgres
        - ms_sql
    - *CONNECTION_STRING* - строка подключения к бд.
    - *DB_TTL* - максимальное время жизни базы данных. Примеры:
        - 1h
        - 30m
        - 3h 15m
    - *JOB_SCHEDULE_CRON* - расписание для джоба удаления бд в cron-формате.
