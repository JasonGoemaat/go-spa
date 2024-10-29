<script>
	import { base } from '$app/paths'
</script>

<nav>
	<a href="{base}/">Home</a>
	<a href="{base}/about">About</a>
	<a href="{base}/test">Test</a>
	<a href="{base}/api/hello">/api/hello</a>
</nav>

<h2>base is: {base}</h2>

<slot></slot>