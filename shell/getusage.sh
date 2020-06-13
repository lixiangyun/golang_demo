#!/bin/bash

GET_CPU_USAGE()
{
    echo "get cpu usage."

    cpu_idle=`top -n 1|grep '%Cpu'|awk '{print $8}'`
    cpu_used=`echo "100 - $cpu_idle"|bc -l`

    if [ `echo $cpu_used|cut -c 1` = "." ]; then
        cpu_used=0$cpu_used
    fi

    echo "cpu_usage: $cpu_used %"

    if [ $(echo "$cpu_used > 90.0" | bc) = 1 ]; then
        echo "abnormal"
        return 1
    else
        echo "normal"
        return 0
    fi
}


GET_MEMORY_USAGE()
{
    echo "get memory usage."
    
    memory_used=`free -m | sed -n '2p'| awk '{print $3}'`
    memory_total=`free -m | sed -n '2p'| awk '{print $2}'`
    memory_rate=$((((memory_used))*100/((memory_total))))
    
    echo "memory_usage: $memory_rate %"
    echo "memory_used : $memory_used MB"
    echo "memory_total: $memory_total MB"

    if [ $memory_rate -gt 90 ]; then
       echo "abnormal"
       return 1 
    else
       echo "normal"
       return 0
    fi
}

GET_DISK_USAGE()
{
    echo "get disk usage."
    
    disk_total=0
    disk_used=0
    disk_max_rate=0
    
    IFS=$'\n'
    for disk in `df | grep -v 'tmpfs' |grep -v 'loop' | grep '/dev/'`; do    
        disk_tmp_name=`echo $disk|awk '{print $1}'`
        disk_tmp_total=`echo $disk|awk '{print $2}'`
        disk_tmp_used=`echo $disk|awk '{print $3}'`
        disk_tmp_rate=`echo $disk|awk '{print $5}'|sed 's/%//g'`
        echo "$disk_tmp_name $disk_tmp_total KB $disk_tmp_used KB $disk_tmp_rate %"
        
        if [ $disk_tmp_rate -gt 90 ]; then
           echo "abnormal"
           return 1 
        fi
        
        if [ $disk_tmp_rate -gt $disk_max_rate ]; then
           disk_max_rate=$((disk_tmp_rate))
        fi
        
        disk_total=$((((disk_total))+((disk_tmp_total))))
        disk_used=$((((disk_used))+((disk_tmp_used))))
    done
    
    disk_rate=$((((disk_used))*100/((disk_total))))
    
    echo "disk_total: $disk_total KB"
    echo "disk_used : $disk_used KB"
    echo "disk_rate : $disk_rate %"
    
    echo "normal"
    return 0
}


GET_CPU_USAGE
GET_MEMORY_USAGE
GET_DISK_USAGE

