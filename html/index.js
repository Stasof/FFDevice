    let cameraStreamUrl="" //будем хранить URL для подключения камеры
        // Логика для прогрессбара
        const progressBar = document.getElementById('progressBar');
        const progressText = document.getElementById('progressText');
               
        const light = document.getElementById('light'); //свет
        const vent_in = document.getElementById('vent_in'); //внктренняя вентиляция
        const vent_ext = document.getElementById('vent_ext'); //внешняя вентиляция

        const nozzleTemp=document.getElementById('nozzleTemp');
        const platTemp=document.getElementById('platTemp');
        
        const rawdetail = document.getElementById('rawdetail') //раздел с информацией detail
        
        const files = document.getElementById('files') //списк файлов в принтере
        const modelthumb = document.getElementById('modelthumb') //превью текущей печатаемой модели
        const status = document.getElementById('status') //статус принтера ready-без задачи, printing-печать, completed-печать завершена
        
        const pauseBtn = document.getElementById('pauseBtn')
        const resumeBtn = document.getElementById('resumeBtn')
        const stopBtn = document.getElementById('stopBtn')

        const remainingDiskSpace = document.getElementById('remainingDiskSpace')

        const refreshBtn = document.getElementById('refreshBtn')

        const dlgPrint=document.getElementById('dlgPrint')
        const printBtn=document.getElementById('printBtn')
        const dlgCancel=document.getElementById('dlgCancel')
        const closeBtn=document.getElementById('closeBtn')

        stopBtn.addEventListener('click', () => {
            dlgCancel.showModal()
        });

        pauseBtn.addEventListener('click', () => {
            Send({cmd:"command", args:"pause"})
        });

        resumeBtn.addEventListener('click', () => {
            Send({cmd:"command", args:"continue"})
        });

        closeBtn.addEventListener("click",()=>{
            Send({cmd:"command", args:"cancel"})
            dlgCancel.close()
        })

        refreshBtn.addEventListener("click",()=>{
            RefreshFiles()
        })

        light.addEventListener("change",()=>{
            Send({cmd:"light", args:light.checked+""})
        })

        vent_in.addEventListener("change",()=>{
            Send({cmd:"fan", args:vent_in.checked+" false"})
            vent_ext.checked=false;
        })
        vent_ext.addEventListener("change",()=>{
            Send({cmd:"fan", args:"false "+vent_ext.checked})
            vent_in.checked=false
        })

        printBtn.addEventListener("click",()=>{
            console.log("PRINT")
            let filename=dlgFileName.innerHTML
            Send({cmd:"print", args:filename})
        })

        async function updateProgress() {
            Send({cmd:"detail", args:""},(ret)=>{
                if(ret)
                {
                    html=""
                    for (const key in ret.detail){
                        html+=`${key}=${ret.detail[key]}<br>`
                    }
                    rawdetail.innerHTML=html
                        
                    progressBar.max= 1 //ret.detail.targetPrintLayer
                    progressBar.value = ret.detail.printProgress//ret.detail.printLayer
                    progressText.textContent = Math.trunc(ret.detail.printProgress*100)+"% "+ret.detail.printFileName;
                    
                    cameraStreamUrl=ret.detail.cameraStreamUrl

                    light.checked=(ret.detail.lightStatus=="open")?true:false;
                    vent_ext.checked=(ret.detail.externalFanStatus=="open")?true:false;
                    vent_in.checked=(ret.detail.internalFanStatus=="open")?true:false;
                    
                    nozzleTemp.innerHTML=Math.trunc(ret.detail.rightTemp)+" / "+ret.detail.rightTargetTemp
                    platTemp.innerHTML=Math.trunc(ret.detail.platTemp)+" / "+ret.detail.platTargetTemp

                    remainingDiskSpace.innerHTML=Math.trunc(ret.detail.remainingDiskSpace*10)/10 + " Гб"
                    cumulativeFilament.innerHTML=Math.trunc(ret.detail.cumulativeFilament*10)/10 + " м"
                    cumulativePrintTime.innerHTML=
                        Math.trunc(ret.detail.cumulativePrintTime/60)+" ч "+ret.detail.cumulativePrintTime%60 + " м"
                    tvoc.innerHTML=ret.detail.tvoc

                    if (ret.detail.printFileThumbUrl && modelthumb.src!=ret.detail.printFileThumbUrl)
                        modelthumb.src=ret.detail.printFileThumbUrl
                    else if(!modelthumb.src.endsWith("printer.png") && !ret.detail.printFileThumbUrl){

                        console.log(modelthumb.src)
                        modelthumb.src="printer.png"
                    }
                    
                    statuses={"ready":"Готов",
                        "printing":"Печать",
                        "cancel":"Отмена",
                        "pause":"Пауза",
                        "heating":"Нагрев",
                        "pausing":"На паузу",
                        "completed":"Завершено",
                        "canceled":"Отменено",
                        "":"Нет соединения"
                    }

                    status.innerHTML=statuses[ret.detail.status]
                    
                    pauseBtn.disabled=true
                    stopBtn.disabled=true
                    resumeBtn.disabled=true
                    if (ret.detail.status=="printing" )
                    {
                        stopBtn.disabled=false
                        pauseBtn.disabled=false
                    }
                    else if (ret.detail.status=="cancel")
                    {
                    }
                    else if (ret.detail.status=="pause" )
                    {
                        resumeBtn.disabled=false
                    }
                    //heating - нагрев после паузы (в это время нельзя поставить на паузу)
                    //pausing - состояние постановки на паузу
                    //completed - печать завершена (требуется подойти к принтеру)
                }
            })  
        }

        // Запуск прогрессбара при загрузке страницы
        progressInterval = setInterval(updateProgress, 3000);
        updateProgress();

        function RefreshFiles()
        {
            Send({cmd:"files", args:""},(ret)=>{
                html="<table>"
                        for (i=0 ;i<ret.gcodeList?.length;i++){
                            let imgid=ret.gcodeList[i]
                            html+=`<tr><td width="64"><img width=64 src="" id="${ret.gcodeList[i]}"></td><td class="selcell" onclick="ShowPrintDialog('${imgid}')">${ret.gcodeList[i]}</td></tr>`
                            Send({Cmd:"thumb",args:imgid},(ret)=>{
                                document.getElementById(imgid).src="data:image/bmp;base64,"+ret.imageData
                            })
                        }
                files.innerHTML=html+"</table>"
            })
        }
        RefreshFiles()

        function ShowPrintDialog(imgid){
            document.getElementById('dlgPreview').src=document.getElementById(imgid).src
            document.getElementById("dlgFileName").innerHTML=imgid
            dlgPrint.showModal()
        }

        function Send(data, callback){
            fetch('http://localhost:8765/api', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
            })
            .then(async response => {
                if (!response.ok) {
                    throw new Error('Ошибка сети или сервера');
                }
                if (callback)
                    callback(await response.json());
            })
            .then(result => {
                console.log('Успех:', result);
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
        }