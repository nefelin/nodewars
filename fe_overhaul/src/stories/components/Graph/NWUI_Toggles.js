function toggle_traffic(context, flag) {
  context.runTrafficDaemon(flag)
}

function toggle_power(context, flag) {
  console.log('draw power')
  context.drawPowerTokens()
}

function toggle_production(context, flag) {
  context.runCoinPulseDaemon(flag)
}


const UI_Toggles = [
  {
    label: 'Traffic',
    icon: './icons/ui_traffic.png',
    handler: toggle_traffic,
  },
  {
    label: 'Power',
    icon: './icons/ui_power3.png',
    sizeMod: 1.2,
    handler: toggle_power,
  },
  {
    label: 'Production',
    icon: './icons/ui_production4.png',
    handler: toggle_production,
  },
  {
    label: 'Alerts',
    icon: './icons/ui_alert.png',
  },
]

export default UI_Toggles