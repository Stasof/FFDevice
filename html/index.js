var FF={}; //Глобальный объект элементов и переменных
(function (){
    document.querySelectorAll('*').forEach(element => {
        if (element.id)
            FF[element.id]=element;
    });
})() 

// обработка заголовков секций
document.addEventListener('DOMContentLoaded', function() {
    const headers = document.querySelectorAll('.section-header'); 
    headers.forEach(header => {
        header.addEventListener('click', function() {
            const section = this.parentElement;
            section.classList.toggle('collapsed');
            if (section.id=="camera")
            { 
                if (!section.classList.contains('collapsed')) {
                    FF.cameraVideo.src = FF.cameraStreamUrl; // Очистка источника
                } else {
                    FF.cameraVideo.src = '';
                }
            }
            else if (section.id=="filelist")
            { 
                if (!section.classList.contains('collapsed')) {
                    FF.refreshBtn.style.display = ''; 
                } else {
                    FF.refreshBtn.style.display = 'none';
                }
            }
            else if (section.id=="information")
            { 
                if (!section.classList.contains('collapsed')) {
                    FF.allinfo.style.display = ''; 
                } else {
                    FF.allinfo.style.display = 'none';
                }
            }
        });
    });
});   

FF.stopBtn.addEventListener('click', () => {FF.dlgCancel.showModal()});
FF.pauseBtn.addEventListener('click', () => {Send({cmd:"command", args:"pause"})});
FF.resumeBtn.addEventListener('click', () => {Send({cmd:"command", args:"continue"})});
FF.closeBtn.addEventListener("click",()=>{
    Send({cmd:"command", args:"cancel"})
    dlgCancel.close()
})
FF.allinfo.addEventListener("click",(event)=>{updateProgress();event.stopPropagation()})
FF.refreshBtn.addEventListener("click",(event)=>{RefreshFiles(); event.stopPropagation()})
FF.light.addEventListener("change",()=>{Send({cmd:"light", args:FF.light.checked+""})})
FF.vent_in.addEventListener("change",()=>{
    Send({cmd:"fan", args:FF.vent_in.checked+" false"})
    FF.vent_ext.checked=false;
})
FF.vent_ext.addEventListener("change",()=>{
    Send({cmd:"fan", args:"false "+FF.vent_ext.checked})
    FF.vent_in.checked=false
})
FF.printBtn.addEventListener("click",()=>{
    let filename=FF.dlgFileName.innerHTML
    Send({cmd:"print", args:JSON.stringify({fileName:filename,levelingBeforePrint:FF.levelingBeforePrint.checked})})
    dlgPrint.close()
})

async function updateProgress() {
    Send({cmd:"detail", args:""},(ret)=>{
        if(ret)
        {
            html="<table>"
            for (const key in ret.detail){
                if (FF.allinfo.checked)
                    html+=`<tr><td>${key}</td><td>${ret.detail[key]}</td></tr>`
                else
                {
                    if (["name","location","firmwareVersion","ipAddr","nozzleModel"].includes(key))
                        html+=`<tr><td>${key}</td><td>${ret.detail[key]}</td></tr>`
                }
            }
            FF.rawdetail.innerHTML=html+"</table>"
                
            FF.progressBar.max= 1 //ret.detail.targetPrintLayer
            FF.progressBar.value = ret.detail.printProgress//ret.detail.printLayer
            FF.progressText.textContent = Math.trunc(ret.detail.printProgress*100)+"% "+ret.detail.printFileName;
            
            FF.cameraStreamUrl=ret.detail.cameraStreamUrl

            FF.light.checked=(ret.detail.lightStatus=="open")?true:false;
            FF.vent_ext.checked=(ret.detail.externalFanStatus=="open")?true:false;
            FF.vent_in.checked=(ret.detail.internalFanStatus=="open")?true:false;
            
            function change(element, newval)
            {
                if (element.innerHTML!=newval)
                    element.parentElement.style="font-weight: bold;"
                else
                    element.parentElement.style="font-weight: normal;"
                return newval
            }

            FF.nozzleTemp.innerHTML=change(FF.nozzleTemp,Math.trunc(ret.detail.rightTemp)+" / "+ret.detail.rightTargetTemp)
            FF.platTemp.innerHTML=change(FF.platTemp,Math.trunc(ret.detail.platTemp)+" / "+ret.detail.platTargetTemp)        

            FF.remainingDiskSpace.innerHTML=Math.trunc(ret.detail.remainingDiskSpace*10)/10 + " Гб"
            FF.cumulativeFilament.innerHTML=Math.trunc(ret.detail.cumulativeFilament*10)/10 + " м"
            FF.cumulativePrintTime.innerHTML=
                Math.trunc(ret.detail.cumulativePrintTime/60)+" ч "+ret.detail.cumulativePrintTime%60 + " м"
            FF.tvoc.innerHTML=change(FF.tvoc,ret.detail.tvoc)

            if (ret.detail.printFileThumbUrl && modelthumb.src!=ret.detail.printFileThumbUrl)
                FF.modelthumb.src=ret.detail.printFileThumbUrl
            else if(!modelthumb.src.endsWith("printer.png") && !ret.detail.printFileThumbUrl)
                FF.modelthumb.src="printer.png"
            
            
            statuses={"ready":"Готов", "printing":"Печать", "cancel":"Отмена", "pause":"Пауза",
                "heating":"Нагрев", "pausing":"На паузу", "completed":"Завершено", "canceled":"Отменено",
                "":"Нет соединения"}
            if (statuses[ret.detail.status]==undefined)
                FF.status.innerHTML=et.detail.status
            else
                FF.status.innerHTML=statuses[ret.detail.status]
            
            FF.pauseBtn.disabled=true
            FF.stopBtn.disabled=true
            FF.resumeBtn.disabled=true
            if (ret.detail.status=="printing" )
            {
                FF.stopBtn.disabled=false
                FF.pauseBtn.disabled=false
            }
            else if (ret.detail.status=="cancel")
            {
            }
            else if (ret.detail.status=="pause" )
            {
                FF.resumeBtn.disabled=false
            }
            //heating - нагрев после паузы (в это время нельзя поставить на паузу)
            //pausing - состояние постановки на паузу
            //completed - печать завершена (требуется подойти к принтеру)
        }
    })  
}

// Запуск прогрессбара при загрузке страницы
setInterval(updateProgress, 3000);
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
        FF.files.innerHTML=html+"</table>"
    })
}
RefreshFiles()

function ShowPrintDialog(imgid){
    FF.dlgPreview.src=document.getElementById(imgid).src
    FF.dlgFileName.innerHTML=imgid
    FF.dlgPrint.showModal()
}

function Send(data, callback){

    // Получаем строку запроса (всё после ?)
    const queryString = window.location.search;
    // Парсим строку запроса
    const urlParams = new URLSearchParams(queryString);
    // Получаем значение конкретного параметра
    const ip = urlParams.get('ip'); 
    const serial = urlParams.get('serial');    
    const check = urlParams.get('check');

    data.printer={ip:ip,serial:serial,check:check}

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
        //console.log('Успех:', result);
    })
    .catch(error => {
        console.error('Ошибка:', error);
    });
}