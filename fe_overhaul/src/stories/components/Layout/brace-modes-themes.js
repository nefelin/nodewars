// console.log('Importing Brace modes/themes...')
// editor modes
import 'brace/mode/c_cpp';
import 'brace/mode/csharp';
import 'brace/mode/clojure';
import 'brace/mode/golang';
import 'brace/mode/haskell';
import 'brace/mode/java';
import 'brace/mode/javascript';
import 'brace/mode/php';
import 'brace/mode/perl';
import 'brace/mode/ruby';
import 'brace/mode/rust';
import 'brace/mode/scala';
import 'brace/mode/sh';
import 'brace/mode/python';

// editor themese
import 'brace/theme/solarized_light';
import 'brace/theme/chrome';
import 'brace/theme/gruvbox';
import 'brace/theme/monokai';

const modeMap = {
	"C++": 'c_cpp',
	"C#": 'csharp',
	"Clojure": 'clojure',
	"Go": 'golang',
	"Haskell": 'haskell',
	"Java": 'java',
	"JavaScript": 'javascript',
	"PHP": 'php',
	"Perl": 'perl',
	"Python": 'python',
	"Ruby": 'ruby',
	"Rust": 'rust',
	"Scala": 'scala',
	"Bash": 'sh', // this is actually sh mode but BE supports Bash and Names are used for lookup
}

const themeMap = {
	'Solarized (light)': 'solarized_light',
	'Chrome    (light)': 'chrome',
	'Gruvbox    (dark)': 'gruvbox',
	'Monokai    (dark)': 'monokai',
}

function modeLookup(aceMode) {
	for (let langName of Object.keys(modeMap)){
		if (modeMap[langName] == aceMode)
			return langName
	}
	return 'error'
}

export {modeMap, themeMap, modeLookup}