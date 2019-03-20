SERVERLIST='routeserver superserver loginserver gatewayserver'

# 服务器关闭顺序
SERVERLIST_REVERSE=
for serv in $SERVERLIST
do
		SERVERLIST_REVERSE=${serv}" "$SERVERLIST_REVERSE
done


DEBUG='true'

dowork()
{
		cp /dev/null nohup.out
		#filename=nohup_`date +"%Y%m%d_%k%M%S"`.log

		case $DEBUG in
				true)

						nohup $PWD/superserver/superserver -logfile=$HOME/log/superserver.log &
						sleep 1

						for ((i=0;i<2;++i))
						do
								nohup $PWD/routeserver/routeserver -logfile=$HOME/log/routeserver$i.log &
								sleep 1
						done
					
						for ((i=0;i<2;++i))
						do
								nohup $PWD/gatewayserver/gatewayserver -logfile=$HOME/log/gatewayserver$i.log -port=802$i &
								sleep 1
						done


						for ((i=0;i<1;++i))
						do
								nohup $PWD/loginserver/loginserver -logfile=$HOME/log/loginserver$i.log -port=801$i &
								sleep 1
						done



						;;
		esac
}

stopwork()
{
		tmpcount=1
		for serv in $SERVERLIST_REVERSE
		do
				if [ ${tmpcount} -eq 1 ]
				then
						echo -n "stoping $serv/$serv "
						tmpcount=$[tmpcount-1]
				fi

				#pkill -9 ${serv:0:10} -u `whoami`
				ps aux|grep "/$serv"|sed -e '/grep/d'|awk '{print $2}'|xargs kill 2&>/dev/null
				while test -f  RunServer.sh
				do  #确保结束第一个进程后再结束第二个，方便MonitorServer监控
						echo -n "."
						COUNT=`ps x|grep -v "grep"|grep "$serv"| wc -l`
						if [ $COUNT -eq 0 ]
						then
								break
						fi
						sleep 0.1
				done
				echo "ok"
				tmpcount=1
		done
}

echo "--------------------------------------------------"
echo "--------------------START-------------------------"
echo "--------------------------------------------------"

case $1 in
		stop)
				stopwork
				;;
		start)
				dowork
				;;
		reboot)
				stopwork
				sleep 3
				dowork
				;;
		dist)
				SERVERDIR='release'
				case $2 in
						stop)
								stopwork
								;;
						start)
								dowork
								;;
						reboot)
								stopwork
								sleep 3
								dowork
								;;
						*)
								stopwork
								sleep 1
								dowork
								;;
				esac
				;;
		*)
				SERVERDIR=
				stopwork
				sleep 1
				dowork
				;;
esac

echo "--------------------------------------------------"
echo "----------------------DONE------------------------"
echo "--------------------------------------------------"


