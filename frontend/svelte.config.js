// import adapter from '@sveltejs/adapter-auto'
import adapter from '@sveltejs/adapter-static'
import RemoteAssets from 'vite-plugin-remote-assets'

/** @type {import('@sveltejs/kit').Config} */
const config = {
    kit: {
    // adapter-auto only supports some environments, see https://kit.svelte.dev/docs/adapter-auto for a list.
    // If your environment is not supported or you settled on a specific environment, switch out the adapter.
    // See https://kit.svelte.dev/docs/adapters for more information about adapters.
        adapter: adapter({
      fallback: '200.html'
    })
  },
  plugins: [
    RemoteAssets()
  ],
  onwarn: (warning, handler) => {
    const IGNORE_WARNINGS = [
      'a11y-click-events-have-key-events',
      'a11y-no-static-element-interactions',
      'a11y-missing-attribute'
    ]
    if (IGNORE_WARNINGS.includes(warning.code)) {
      return
        }
    handler(warning)
    }
}

export default config
