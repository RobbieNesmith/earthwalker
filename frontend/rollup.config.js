import svelte from 'rollup-plugin-svelte';
import { nodeResolve } from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import livereload from 'rollup-plugin-livereload';
import { terser } from 'rollup-plugin-terser';
import css from 'rollup-plugin-css-only';
import fs from "fs";
import posthtml from "posthtml";
import { hash } from "posthtml-hash";
import copy from "rollup-plugin-copy";
import nodePolyfills from 'rollup-plugin-node-polyfills';


const production = !process.env.ROLLUP_WATCH;
const BUILD_DIR = '../public';

// Inspired by https://github.com/metonym/svelte-rollup-template/blob/master/rollup.config.js
function hashStatic(htmlPath) {
	return {
	  name: "hash-static",
	//   buildStart() {
	// 	rimraf.sync(OUT_DIR);
	//   },
	  writeBundle() {
		posthtml().use(
		  // hashes `bundle.[custom-hash].css` and `bundle.[custom-hash].js`
		  hash({ path: BUILD_DIR, pattern: new RegExp(/\[custom-hash\]/), }),
		)
		  .process(fs.readFileSync(htmlPath))
		  .then((result) =>
			fs.writeFileSync(htmlPath, result.html)
		  );
	  },
	};
  }

export default {
	input: [
		'src/js/main.js',
	],
	output: {
		sourcemap: 'inline',
		format: 'iife',
		name: 'app',
		dir: BUILD_DIR,
		entryFileNames: '[name]-[custom-hash].js',
	},
	plugins: [
		copy({ targets: [{ src: "src/**/*.html", dest: BUILD_DIR }] }),
		svelte({
			compilerOptions: {
				// enable run-time checks when not in production
				dev: !production,
			}
		}),

		// we'll extract any component CSS out into
		// a separate file - better for performance
		css({
			output: 'bundle-[custom-hash].css',

		}),

		// If you have external dependencies installed from
		// npm, you'll most likely need these plugins. In
		// some cases you'll need additional configuration -
		// consult the documentation for details:
		// https://github.com/rollup/plugins/tree/master/packages/commonjs
		nodeResolve({
			browser: true,
			dedupe: ['svelte'],
		}),
		commonjs(),



		// In dev mode, call `npm run start` once
		// the bundle has been generated
		!production && serve(),

		// Watch the `public` directory and refresh the
		// browser on changes when not in production
		!production && livereload(BUILD_DIR),

		// If we're building for production (npm run build
		// instead of npm run dev), minify
		production && terser(),
		production && hashStatic(`${BUILD_DIR}/index.html`),
	],
	watch: {
		clearScreen: false
	}
};

function serve() {
	let started = false;

	return {
		writeBundle() {
			if (!started) {
				started = true;

				require('child_process').spawn('npm', ['run', 'start', '--', '--dev'], {
					stdio: ['ignore', 'inherit', 'inherit'],
					shell: true
				});
			}
		}
	};
}
