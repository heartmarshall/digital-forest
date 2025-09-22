import os

def collect_files_to_txt(directory_path, output_file_name):
    """
    Собирает содержимое всех файлов из указанной директории в один TXT-файл.

    :param directory_path: Путь к директории, файлы из которой нужно собрать.
    :param output_file_name: Имя итогового TXT-файла.
    """
    try:
        # Открываем или создаем файл для записи. 'w' означает, что файл будет перезаписан,
        # если он уже существует. Используйте 'a' для добавления в конец файла.
        # Указываем кодировку utf-8 для поддержки различных символов.
        with open(output_file_name, 'w', encoding='utf-8') as output_file:
            # os.walk() рекурсивно обходит все директории и файлы.
            for root, dirs, files in os.walk(directory_path):
                for file_name in files:
                    # Формируем полный путь к файлу.
                    file_path = os.path.join(root, file_name)
                    
                    try:
                        # Открываем и читаем содержимое текущего файла.
                        with open(file_path, 'r', encoding='utf-8', errors='ignore') as current_file:
                            content = current_file.read()
                        
                        # Получаем абсолютный путь к файлу.
                        full_path = os.path.abspath(file_path)

                        # Записываем информацию в итоговый файл в нужном формате.
                        output_file.write(f"Название файла - {file_name}\n")
                        output_file.write(f"полный путь до него - {full_path}\n")
                        output_file.write("содержимое\n")
                        output_file.write(content)
                        output_file.write("\n\n" + "="*80 + "\n\n")

                    except Exception as e:
                        print(f"Не удалось прочитать файл: {file_path}. Ошибка: {e}")

        print(f"Все файлы из директории '{directory_path}' были успешно собраны в файл '{output_file_name}'.")

    except FileNotFoundError:
        print(f"Ошибка: Указанная директория не найдена '{directory_path}'.")
    except Exception as e:
        print(f"Произошла непредвиденная ошибка: {e}")

# --- Использование скрипта ---

# 1. Укажите путь к директории, которую нужно обработать.
#    Пример для Windows: 'C:\\Users\\ИмяПользователя\\Documents'
#    Пример для macOS/Linux: '/home/ИмяПользователя/Documents'
target_directory = '/home/alodi/playground/digital-forest/digital-forest-frontend copy/' 

# 2. Укажите имя файла, в который будут собраны все данные.
output_file = 'собранные_файлы2.txt'

# 3. Запуск функции.
if target_directory != 'путь_к_вашей_директории':
    collect_files_to_txt(target_directory, output_file)
else:
    print("Пожалуйста, укажите путь к директории в переменной 'target_directory'.")