package frontendmisc

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"
)

// RenderPackage renders the package
func RenderPackage(project *models.Project) error {
	files := map[string]func() []byte{
		"frontend/v1/README.md":        readmeDotMD,
		"frontend/v1/rollup.config.js": rollupDotConfigDotJS,
		"frontend/v1/package.json":     packageDotJSON,
	}

	for filename, file := range files {
		fname := utils.BuildTemplatePath(project.OutputPath, filename)

		if mkdirErr := os.MkdirAll(filepath.Dir(fname), os.ModePerm); mkdirErr != nil {
			return fmt.Errorf("error making directory: %v\n", mkdirErr)
		}

		f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		bytes := file()
		if _, err := f.Write(bytes); err != nil {
			return err
		}
	}

	return nil
}

func readmeDotMD() []byte {
	return []byte(`# Get started

Install the dependencies...

` + "```" + `bash
cd svelte-app
npm install
` + "```" + `

...then start [Rollup](https://rollupjs.org):

` + "```" + `bash
npm run dev
` + "```" + `

Navigate to [localhost:5000](http://localhost:5000). You should see your app running. Edit a component file in ` + "`" + `src` + "`" + `, save it, and reload the page to see your changes.
`)
}

func rollupDotConfigDotJS() []byte {
	return []byte(`import svelte from 'rollup-plugin-svelte';
import resolve from 'rollup-plugin-node-resolve';
import commonjs from 'rollup-plugin-commonjs';
import livereload from 'rollup-plugin-livereload';
import { terser } from 'rollup-plugin-terser';

const production = !process.env.ROLLUP_WATCH;

export default {
	input: 'src/main.js',
	output: {
		sourcemap: true,
		format: 'iife',
		name: 'app',
		file: 'public/bundle.js'
	},
	plugins: [
		svelte({
			// enable run-time checks when not in production
			dev: !production,
			// we'll extract any component CSS out into
			// a separate file — better for performance
			css: css => {
				css.write('public/bundle.css');
			}
		}),

		// If you have external dependencies installed from
		// npm, you'll most likely need these plugins. In
		// some cases you'll need additional configuration —
		// consult the documentation for details:
		// https://github.com/rollup/rollup-plugin-commonjs
		resolve(),
		commonjs(),

		// Watch the ` + "`" + `public` + "`" + ` directory and refresh the
		// browser on changes when not in production
		!production && livereload('public'),

		// If we're building for production (npm run build
		// instead of npm run dev), minify
		production && terser()
	],
	watch: {
		clearScreen: false
	}
};
`)
}

func packageDotJSON() []byte {
	return []byte(`{
	"name": "svelte-app",
	"version": "1.0.0",
	"devDependencies": {
		"npm-run-all": "^4.1.5",
		"rollup": "^1.10.1",
		"rollup-plugin-commonjs": "^9.3.4",
		"rollup-plugin-livereload": "^1.0.4",
		"rollup-plugin-node-resolve": "^4.2.3",
		"rollup-plugin-svelte": "^5.0.3",
		"rollup-plugin-terser": "^4.0.4",
		"sirv-cli": "^0.4.0",
		"svelte": "^3.0.0"
	},
	"scripts": {
		"build": "rollup -c",
		"autobuild": "rollup -c -w",
		"dev": "run-p start:dev autobuild",
		"start": "sirv public",
		"start:dev": "sirv public --dev"
	},
	"dependencies": {
		"svelte-routing": "^1.1.1"
	}
}`)

}
