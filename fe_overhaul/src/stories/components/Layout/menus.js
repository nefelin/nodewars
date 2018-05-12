export const Menus = {
	Main: [{ name: 'NodeWars' }, {name: 'Layout', items: ['Load', 'Save', 'Reset']}],
	Editor: [{ name: 'Ace Editor' }, { name: 'Build', items: ['Make, (ctrl-enter)', 'Test, (ctrl-\\)', 'Reset, (ctrl-r)']}, { name: 'Theme' }, { name: 'Language' }],
	Score: [{ name: 'Score' }],
	MapMenu: [{ name: 'Map' }],
	Terminal: [{ name: 'Terminal' }],

	Challenge: [{ name: 'Challenge Details' }],
	Results: [{ name: 'Test Results' }],
	CompilerOutput: [{ name: 'Compiler Output' }],
	Stdin: [{ name: 'Stdin' }],
}

export function buildEditorMenu(modeMap, themeMap, state) {
	return ([
		{ name: 'Ace Editor' },
		{ name: 'Build', items: ['Make, (ctrl-enter)', 'Test, (ctrl-\\)', 'Reset, (ctrl-r)']},
		{ name: 'Theme', items: Object.keys(themeMap) },
		{ name: 'Language: ' + state.aceMode, items: Object.keys(modeMap) }
	])
}




//Editor: [{ name: 'Ace Editor' }, { name: 'Build', items: ['Make, (ctrl-enter)', 'Test, (ctrl-\\)', 'Reset, (ctrl-r)']}, { name: 'Theme', items: Object.keys(themeMap) }, { name: 'Language: ' + this.state.aceMode, items: Object.keys(modeMap) }],