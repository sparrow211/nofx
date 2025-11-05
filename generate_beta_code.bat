@echo off
setlocal enabledelayedexpansion

REM ================================================================
REM 内测码生成脚本 - Windows 批处理版本
REM ================================================================
REM 
REM 功能说明:
REM   生成6位不重复的内测码并保存到文件
REM   支持批量生成、列表查看、自定义文件等功能
REM 
REM 使用方法:
REM   generate_beta_code.bat [选项]
REM 
REM 参数选项:
REM   -c COUNT    生成内测码数量 (默认: 1)
REM   -l          列出现有内测码
REM   -f FILE     指定内测码文件路径 (默认: beta_codes.txt)
REM   -h          显示帮助信息
REM 
REM 使用示例:
REM   generate_beta_code.bat              # 生成1个内测码
REM   generate_beta_code.bat -c 10        # 生成10个内测码
REM   generate_beta_code.bat -l           # 列出现有内测码
REM   generate_beta_code.bat -f custom.txt -c 5  # 在自定义文件生成5个
REM 
REM 内测码规则:
REM   - 长度: 6位
REM   - 字符集: 数字2-9, 小写字母a-z
REM   - 排除易混淆字符: 0, 1, i, l, o
REM   - 每个内测码唯一且不重复
REM 
REM ================================================================

REM 设置UTF-8编码支持中文显示
chcp 65001 > nul

REM 显示程序标题
echo.
echo ================================================================
echo                    内测码生成工具 v1.0
echo ================================================================
echo.

REM 初始化默认参数
set "BETA_CODES_FILE=beta_codes.txt"
set "COUNT=1"
set "LIST_ONLY=0"
set "SHOW_HELP=0"
set "CODE_LENGTH=6"
set "INTERACTIVE_MODE=0"

REM 字符集定义 (避免易混淆字符：0/O, 1/I/l)
set "CHARSET=23456789abcdefghjkmnpqrstuvwxyz"

REM 检查是否有命令行参数，如果没有则进入交互模式
if "%~1"=="" (
    set "INTERACTIVE_MODE=1"
    goto :interactive_menu
)

REM ================================================================
REM 解析命令行参数
REM ================================================================
:parse_args
if "%~1"=="" goto :args_done

if /i "%~1"=="-c" (
    set "COUNT=%~2"
    shift
    shift
    goto :parse_args
)

if /i "%~1"=="-l" (
    set "LIST_ONLY=1"
    shift
    goto :parse_args
)

if /i "%~1"=="-f" (
    set "BETA_CODES_FILE=%~2"
    shift
    shift
    goto :parse_args
)

if /i "%~1"=="-h" (
    set "SHOW_HELP=1"
    shift
    goto :parse_args
)

if /i "%~1"=="--help" (
    set "SHOW_HELP=1"
    shift
    goto :parse_args
)

REM 无效参数处理
echo 错误: 无效参数 "%~1"
echo 使用 -h 查看帮助信息
echo.
echo 按任意键退出...
pause > nul
exit /b 1

:args_done

REM ================================================================
REM 交互式菜单
REM ================================================================
:interactive_menu
if %INTERACTIVE_MODE%==0 goto :skip_interactive

REM 清屏并显示标题
cls
echo.
echo ================================================================
echo                    内测码生成工具 v1.0
echo ================================================================

echo ================================================================
echo                        功能菜单  
echo ================================================================
echo.
echo 当前状态:
REM 显示当前文件信息
if exist "%BETA_CODES_FILE%" (
    set "file_count=0"
    for /f "usebackq delims=" %%a in ("%BETA_CODES_FILE%") do (
        set "line=%%a"
        if not "!line!"=="" (
            echo !line! | findstr /r "^#" > nul
            if errorlevel 1 set /a file_count+=1
        )
    )
    echo   文件: %BETA_CODES_FILE% ^(!file_count! 个内测码^)
) else (
    echo   文件: %BETA_CODES_FILE% ^(不存在^)
)
echo.
echo 快速生成:
echo   [1] 生成 1 个内测码
echo   [2] 生成 5 个内测码  
echo   [3] 生成 10 个内测码
echo.
echo 高级选项:
echo   [4] 生成自定义数量内测码
echo   [5] 列出现有内测码
echo   [6] 自定义文件路径生成
echo.
echo 其他:
echo   [7] 查看帮助信息
echo   [8] 退出程序
echo.
echo ================================================================

:menu_choice
set /p "choice=请输入选项 (1-8): "

if "%choice%"=="1" (
    set "COUNT=1"
    set "LIST_ONLY=0"
    echo.
    echo 正在准备生成 1 个内测码...
    goto :skip_interactive
)

if "%choice%"=="2" (
    set "COUNT=5" 
    set "LIST_ONLY=0"
    echo.
    echo 正在准备生成 5 个内测码...
    goto :skip_interactive
)

if "%choice%"=="3" (
    set "COUNT=10"
    set "LIST_ONLY=0"
    echo.
    echo 正在准备生成 10 个内测码...
    goto :skip_interactive
)

if "%choice%"=="4" (
    echo.
    set /p "COUNT=请输入要生成的内测码数量: "
    echo %COUNT% | findstr /r "^[1-9][0-9]*$" > nul
    if errorlevel 1 (
        echo 错误: 请输入有效的正整数
        echo.
        goto :menu_choice
    )
    set "LIST_ONLY=0"
    goto :skip_interactive
)

if "%choice%"=="5" (
    set "LIST_ONLY=1"
    goto :skip_interactive
)

if "%choice%"=="6" (
    echo.
    set /p "BETA_CODES_FILE=请输入文件路径 (默认: beta_codes.txt): "
    if "%BETA_CODES_FILE%"=="" set "BETA_CODES_FILE=beta_codes.txt"
    echo.
    set /p "COUNT=请输入要生成的内测码数量: "
    echo %COUNT% | findstr /r "^[1-9][0-9]*$" > nul
    if errorlevel 1 (
        echo 错误: 请输入有效的正整数
        echo.
        goto :menu_choice
    )
    set "LIST_ONLY=0"
    goto :skip_interactive
)

if "%choice%"=="7" (
    set "SHOW_HELP=1"
    goto :skip_interactive
)

if "%choice%"=="8" (
    echo.
    echo 谢谢使用内测码生成工具!
    echo.
    echo 按任意键退出...
    pause > nul
    exit /b 0
)

echo.
echo 无效选择，请输入 1-8 之间的数字
echo.
goto :menu_choice

:skip_interactive

REM ================================================================
REM 显示帮助信息
REM ================================================================
if %SHOW_HELP%==1 (
    echo.
    echo ================================================================
    echo                    内测码生成工具 - 帮助信息
    echo ================================================================
    echo.
    echo 用法: %~nx0 [选项]
    echo.
    echo 选项:
    echo   -c COUNT    生成内测码数量 ^(默认: 1^)
    echo   -l          列出现有内测码
    echo   -f FILE     内测码文件路径 ^(默认: beta_codes.txt^)
    echo   -h          显示此帮助信息
    echo.
    echo 示例:
    echo   %~nx0                          # 生成1个内测码
    echo   %~nx0 -c 10                    # 生成10个内测码
    echo   %~nx0 -l                       # 列出现有内测码
    echo   %~nx0 -f custom.txt -c 5       # 在自定义文件中生成5个内测码
    echo.
    echo 内测码规则:
    echo   - 长度: %CODE_LENGTH% 位
    echo   - 字符集: 数字 2-9, 小写字母 a-z
    echo   - 排除字符: 0, 1, i, l, o ^(避免视觉混淆^)
    echo   - 每个内测码唯一且不重复
    echo.
    echo ================================================================
    echo.
    
    REM 根据模式决定下一步操作  
    if %INTERACTIVE_MODE%==1 (
        echo [1] 返回主菜单
        echo [2] 退出程序
        echo.
        set /p "next_choice=请选择 (1-2): "
        if "!next_choice!"=="1" (
            echo.
            goto :interactive_menu
        )
    )
    
    echo 按任意键退出...
    pause > nul
    exit /b 0
)

REM ================================================================
REM 验证COUNT参数 (仅在非交互模式下验证，交互模式已经验证过)
REM ================================================================
if %INTERACTIVE_MODE%==0 (
    echo %COUNT% | findstr /r "^[1-9][0-9]*$" > nul
    if errorlevel 1 (
        echo 错误: count 必须是正整数，当前值: %COUNT%
        echo.
        echo 按任意键退出...
        pause > nul
        exit /b 1
    )
)

REM ================================================================
REM 列出现有内测码功能
REM ================================================================
if %LIST_ONLY%==1 (
    echo.
    echo 正在读取内测码文件: %BETA_CODES_FILE%
    echo.
    
    if exist "%BETA_CODES_FILE%" (
        set "code_count=0"
        set "line_num=0"
        
        REM 统计有效内测码数量
        for /f "usebackq delims=" %%a in ("%BETA_CODES_FILE%") do (
            set "line=%%a"
            REM 跳过空行和注释行
            if not "!line!"=="" (
                echo !line! | findstr /r "^#" > nul
                if errorlevel 1 (
                    set /a code_count+=1
                )
            )
        )
        
        if !code_count!==0 (
            echo   内测码列表为空
        ) else (
            echo 当前内测码 ^(!code_count! 个^):
            echo.
            REM 显示内测码列表
            for /f "usebackq delims=" %%a in ("%BETA_CODES_FILE%") do (
                set "line=%%a"
                if not "!line!"=="" (
                    echo !line! | findstr /r "^#" > nul
                    if errorlevel 1 (
                        set /a line_num+=1
                        echo   !line_num!. %%a
                    )
                )
            )
        )
        echo.
        echo 总计: !code_count! 个内测码
    ) else (
        echo   内测码文件不存在: %BETA_CODES_FILE%
        echo   提示: 运行生成命令创建新文件
    )
    echo.
    
    REM 根据模式决定下一步操作
    if %INTERACTIVE_MODE%==1 (
        echo [1] 返回主菜单
        echo [2] 退出程序
        echo.
        set /p "next_choice=请选择 (1-2): "
        if "!next_choice!"=="1" (
            echo.
            goto :interactive_menu
        )
    )
    
    echo 按任意键退出...
    pause > nul
    exit /b 0
)

REM ================================================================
REM 生成内测码主逻辑
REM ================================================================

echo.
echo ================================================================
echo 开始生成内测码
echo ================================================================
echo 数量: %COUNT% 个
echo 文件: %BETA_CODES_FILE%
echo 模式: 交互模式=%INTERACTIVE_MODE%
echo ================================================================
echo.

REM 计算字符集长度
set "charset_len=0"
:count_charset
call set "char=%%CHARSET:~%charset_len%,1%%"
if "%char%"=="" goto :charset_counted
set /a charset_len+=1
goto :count_charset
:charset_counted

echo 字符集长度: %charset_len%
echo 开始生成...

REM 验证字符集长度
if %charset_len% leq 0 (
    echo 错误: 字符集长度计算失败
    echo.
    echo 按任意键退出...
    pause > nul
    exit /b 1
)

REM 生成内测码
set "generated_count=0"
set "new_codes="
set "max_attempts=1000"

:generate_loop
if %generated_count% geq %COUNT% goto :generate_done

set "attempts=0"

:retry_generate
REM 检查最大尝试次数
if %attempts% geq %max_attempts% (
    echo 警告: 生成第 !generated_count! 个内测码时达到最大尝试次数
    echo 提示: 可能字符空间不足，请尝试减少生成数量
    goto :generate_done
)

REM 生成单个内测码
set "code="
for /l %%i in (1,1,%CODE_LENGTH%) do (
    set /a "rand_idx=!random! %% %charset_len%"
    call set "code=!code!%%CHARSET:~!rand_idx!,1%%"
)

REM 验证内测码格式 (检查长度)
if not "!code:~6,1!"=="" (
    set /a attempts+=1
    goto :retry_generate
)
if "!code:~5,1!"=="" (
    set /a attempts+=1
    goto :retry_generate
)

REM 检查是否已存在于文件中
if exist "%BETA_CODES_FILE%" (
    for /f "usebackq delims=" %%a in ("%BETA_CODES_FILE%") do (
        if "%%a"=="!code!" (
            set /a attempts+=1
            goto :retry_generate
        )
    )
)

REM 检查是否与本次生成的重复
echo !new_codes! | findstr /c:"!code!" > nul
if not errorlevel 1 (
    set /a attempts+=1
    goto :retry_generate
)

REM 验证字符集 (检查是否包含非法字符)
set "valid_code=1"
for /l %%i in (0,1,5) do (
    set "char=!code:~%%i,1!"
    echo %CHARSET% | findstr /c:"!char!" > nul
    if errorlevel 1 (
        set "valid_code=0"
        goto :check_valid_done
    )
)
:check_valid_done

if %valid_code%==0 (
    set /a attempts+=1
    goto :retry_generate
)

REM 添加到新生成列表和文件
set "new_codes=!new_codes! !code!"
echo !code! >> "%BETA_CODES_FILE%"
set /a generated_count+=1

REM 显示进度
if !generated_count! leq 10 (
    echo   生成内测码 !generated_count!: !code!
) else (
    REM 大量生成时减少输出
    set /a "show_progress=!generated_count! %% 10"
    if !show_progress!==0 (
        echo   已生成 !generated_count! 个内测码...
    )
)

goto :generate_loop

:generate_done

REM ================================================================
REM 处理生成结果
REM ================================================================

if %generated_count%==0 (
    echo.
    echo 错误: 未能生成任何新的内测码
    echo 可能原因: 字符空间不足或文件权限问题
    echo.
    echo 按任意键退出...
    pause > nul
    exit /b 1
)

REM 去重并排序内测码文件
if exist "%BETA_CODES_FILE%" (
    echo.
    echo 正在去重并排序内测码文件...
    
    REM 创建临时文件进行去重排序
    set "temp_file=%BETA_CODES_FILE%.tmp"
    if exist "!temp_file!" del "!temp_file!"
    
    REM 读取所有非空非注释行到临时文件
    for /f "usebackq delims=" %%a in ("%BETA_CODES_FILE%") do (
        set "line=%%a"
        if not "!line!"=="" (
            echo !line! | findstr /r "^#" > nul
            if errorlevel 1 (
                echo %%a >> "!temp_file!"
            )
        )
    )
    
    REM 排序去重 (简单实现)
    if exist "!temp_file!" (
        sort "!temp_file!" > "%BETA_CODES_FILE%"
        del "!temp_file!"
    )
)

REM ================================================================
REM 显示生成结果
REM ================================================================

echo.
echo ================================================================
echo                        生成结果
echo ================================================================
echo.
echo 成功生成 %generated_count% 个内测码:

REM 显示新生成的内测码
set "display_count=0"
for %%c in (%new_codes%) do (
    if not "%%c"=="" (
        set /a display_count+=1
        echo   !display_count!. %%c
    )
)

echo.
echo 内测码文件: %BETA_CODES_FILE%

REM 统计文件中总的内测码数量
if exist "%BETA_CODES_FILE%" (
    set "total_count=0"
    for /f "usebackq delims=" %%a in ("%BETA_CODES_FILE%") do (
        set "line=%%a"
        if not "!line!"=="" (
            echo !line! | findstr /r "^#" > nul
            if errorlevel 1 (
                set /a total_count+=1
            )
        )
    )
    echo 当前内测码总计: !total_count! 个
)

REM 显示内测码规则说明
echo.
echo 内测码规则:
echo   - 长度: %CODE_LENGTH% 位字符
echo   - 字符集: 数字 2-9, 小写字母 a-z
echo   - 排除字符: 0, 1, i, l, o ^(避免视觉混淆^)
echo   - 每个内测码唯一且不重复
echo   - 文件编码: UTF-8

echo.
echo ================================================================
echo 生成完成! 使用 -l 参数可查看所有内测码
echo ================================================================

REM 根据模式决定下一步操作
if %INTERACTIVE_MODE%==1 (
    echo.
    echo [1] 返回主菜单
    echo [2] 退出程序
    echo.
    set /p "next_choice=请选择 (1-2): "
    if "!next_choice!"=="1" (
        echo.
        goto :interactive_menu
    )
)

echo.
echo 按任意键退出...
pause > nul

endlocal
exit /b 0