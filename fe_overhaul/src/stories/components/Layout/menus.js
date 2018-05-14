export const Menus = {
	Main: [{ name: 'NodeWars' }, {name: 'Layout', items: ['Load', 'Save', 'Reset']}],
	// AceMenu: [{ name: 'Ace Editor' }, { name: 'Build', items: ['Make, (ctrl-enter)', 'Test, (ctrl-\\)', 'Reset, (ctrl-r)']}, { name: 'Theme' }],
	Score: [{ name: 'Score' }],
	MapMenu: [{ name: 'Map' }],
	Terminal: [{ name: 'Terminal' }],

	Challenge: [{ name: 'Challenge Details' }],
	Results: [{ name: 'Test Results' }],
	CompilerOutput: [{ name: 'Compiler Output' }],
	Stdin: [{ name: 'Stdin' }],
}

export function AceMenu(modeMap, themeMap, langList, state) {
	const languages = []
	console.log('acemenu state', this.state)
	const menu = [
		{ name: 'Ace Editor' },
		{ name: 'Theme', items: Object.keys(themeMap) },	
	]

	for (let lang of langList) {
		for (let langLabel of Object.keys(modeMap)){
			// console.log('comparing',lang,'to',langLabel)
			if (lang == langLabel.toLowerCase()){
				languages.push(langLabel)
				// console.log('Add language', langLabel)
			}
		}
	}
	
	console.log('Supported Languages:', languages)
	if (languages.length > 0) {
		menu.push({ name: 'Build', items: ['Make, (ctrl-enter)', 'Test, (ctrl-\\)', 'Reset, (ctrl-r)']})
		menu.push({ name: 'Language: ' + state.aceMode, items: languages })
	}

	return menu
}




//Editor: [{ name: 'Ace Editor' }, { name: 'Build', items: ['Make, (ctrl-enter)', 'Test, (ctrl-\\)', 'Reset, (ctrl-r)']}, { name: 'Theme', items: Object.keys(themeMap) }, { name: 'Language: ' + this.state.aceMode, items: Object.keys(modeMap) }],