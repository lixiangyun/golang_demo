#!/bin/bash


GET_TOP_CPU()
{
    # 根据需要进行修改
    TOPNUM=5
    
    OUTPUT=""
    INDEX=-1
    IFS=$'\n'
    for cpu in `top -s -bn 1 |grep -A $((TOPNUM)) "COMMAND"`; do
        INDEX=$(($INDEX+1))
        if [ $INDEX -eq 0 ]; then
            continue
        fi
        cpu_pid=`echo $cpu|awk '{print $1}'`
        cpu_usage=`echo $cpu|awk '{print $9}'`
        cpu_commond=`echo $cpu|awk '{print $12}'`
        
        if [ $(echo "$cpu_usage > 0.0" | bc) = 1 ]; then
            TEMP=`echo "TOP$INDEX: PID=$cpu_pid CPU=$cpu_usage% CMD=$cpu_commond \n"`
            OUTPUT=${OUTPUT}${TEMP}
        fi
        if [ $INDEX -gt $TOPNUM ]; then
            break
        fi
    done
    
    echo -e $OUTPUT
}

GET_TOP_MEM()
{
    TOP_NUM=5
    declare -A topNMap
    declare -A allMap
    
    OUTPUT=""

    INDEX=-1
    IFS=$'\n'
    for mem in `top -bn 1 |grep -A 65535 "COMMAND"`;
    do
        INDEX=$(($INDEX+1))
        if [ $INDEX -eq 0 ];then
            continue
        fi
        
        MEM_PID=`echo $mem | awk '{print $1}'`
        MEM_USAGE=`echo $mem | awk '{print $10}'`
        MEM_CMD=`echo $mem | awk '{print $12}'`
        
        if [ `echo $MEM_USAGE == 0.0 | bc` = 1 ];then
            continue
        fi
        
        allMap[$MEM_PID]=`echo $MEM_PID $MEM_USAGE $MEM_CMD`
    done
    
    # 遍历获取TOPN信息
    for((i=1;i<=$TOP_NUM;i++));do
    
        MAX_USAGE=0.0
        MAX_VALUE=""
        MAX_KEY=""
        for key in ${!allMap[@]};do
            
            # 判断是否已经记录
            value=${topNMap[$key]}

            length=$(echo $value|wc -c)
            if [ $length -gt 1 ];then
                continue
            fi
            
            value=${allMap[$key]}
            temp_usage=`echo $value | awk '{print $2}'`
            
            result=$(echo "$temp_usage>$MAX_USAGE"|bc)
            if [ $result -eq 1 ];then
                MAX_VALUE=$value
                MAX_KEY=$key
                MAX_USAGE=$temp_usage
            fi
        done
        
        length=$(echo $MAX_KEY|wc -c)
        if [ $length -eq 1 ];then
            continue
        fi
        
        topNMap[$MAX_KEY]=$MAX_VALUE
        
        MEM_PID=`echo $MAX_VALUE | awk '{print $1}'`
        MEM_USAGE=`echo $MAX_VALUE | awk '{print $2}'`
        MEM_CMD=`echo $MAX_VALUE | awk '{print $3}'`
        
        OUTPUT=$OUTPUT`echo "PID:$MEM_PID MEM:$MEM_USAGE% CMD:$MEM_CMD\n"`
        
    done
    
    echo -e $OUTPUT
}


GET_TOP_DISK()
{
    TOP_NUM=5
    declare -A topNMap
    declare -A allMap
    
    OUTPUT=""

    IFS=$'\n'
    for disk in `df |grep -v 'tmpfs'|grep -v 'loop'|grep '/dev/'`;
    do
        DISK_NAME=`echo $disk | awk '{print $1}'`
        DISK_USAGE=`echo $disk | awk '{print $5}' | sed 's/%//g'`
        DISK_MOUNT=`echo $disk | awk '{print $6}'`
        
        allMap[$DISK_NAME]=`echo $DISK_NAME $DISK_USAGE $DISK_MOUNT`
    done
    
    # 遍历获取TOPN信息
    for((i=1;i<=$TOP_NUM;i++));do
    
        MAX_USAGE=0
        MAX_VALUE=""
        MAX_KEY=""
        
        for key in ${!allMap[@]};do
            
            # 判断是否已经记录
            value=${topNMap[$key]}

            length=$(echo $value|wc -c)
            if [ $length -gt 1 ];then
                continue
            fi
            
            value=${allMap[$key]}
            temp_usage=`echo $value | awk '{print $2}'`
            
            result=$(echo "$temp_usage>$MAX_USAGE"|bc)
            if [ $result -eq 1 ];then
                MAX_VALUE=$value
                MAX_KEY=$key
                MAX_USAGE=$temp_usage
            fi
        done
        
        length=$(echo $MAX_KEY|wc -c)
        if [ $length -eq 1 ];then
            continue
        fi
        
        topNMap[$MAX_KEY]=$MAX_VALUE
        
        DISK_NAME=`echo $MAX_VALUE | awk '{print $1}'`
        DISK_USAGE=`echo $MAX_VALUE | awk '{print $2}'`
        DISK_MOUNT=`echo $MAX_VALUE | awk '{print $3}'`
        
        OUTPUT=$OUTPUT`echo "FileSystem:$DISK_NAME Use:$DISK_USAGE% Mounted:$DISK_MOUNT\n"`
        
    done
    
    echo -e $OUTPUT
}

GET_TOP_MEM
GET_TOP_DISK
GET_TOP_CPU
