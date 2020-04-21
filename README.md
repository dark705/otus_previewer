# Превьювер изображений
Сервис представляет собой web-сервер (прокси), загружающий изображения, масштабирующий/обрезающий их до нужного формата 
и возвращающий пользователю.
Наиболее часто использующиеся изображения остаются в кэше, редко используемые вытесняются новыми.

####Конфигурация

Настройка осуществляется на основании ENV переменных окружения, с предопределёнными дефолтными значениями:

    LOG_LEVEL - уровень логирования ("error", "warn", "info", "debug"), по умолчанию: "debug"
    HTTP_LISTEN - адрес:порт, на котором запущен сервис, по умолчанию: ":8013"
    IMAGE_MAX_FILE_SIZE - максимальный размер запрашиваемого (исходного) изображения в байтах, по умолчанию: "1000000" (1M)
    CACHE_SIZE - общий размер кэша для всех обработанных и изменённых картинок в байтах по умолчанию: "100000000" (100M)
    CACHE_TYPE - тип кэша ("inmemory" - в оперативной памяти, "disk" - указанная папка на диске), по умолчанию: "disk"
    CACHE_TYPE - путь к папке кэша на диске, по умолчанию "./cache"
    
####Параметры URL запроса

    http://address:port/service/width/height/somesite.com/image.jpg
    
где: 

    address:port - адрес:порт, где запущен сервер
    service - тип операции над изображением 
        - "resize" - изменить размер, пропорции не сохраняются. 
        - "fill" - заполнить, пропорции сохраняются, исходное изображение при этом центрируется и может быть обрезано по высоте или ширине
        - "fit" - вписать изображение в заданный размер, пропорции сохраняются, но полученное изображение может быть меньше по высоте или ширине
    width - ширина
    height - высота 
    somesite.com/image.jpg - адрес до изображения на стороннем ресурсе
    
Выбор тип протокола (https, http), происходит автоматически, и его не надо указывать в пути до удалённого изображения. 
С начала пытается получить данные с удалённого сервера по https, затем в случае ошибки по http.

####Коммиляция
    make