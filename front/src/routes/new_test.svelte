<script>
	import { post, getUser } from './utils';
	import { onMount } from 'svelte';
	import { goto } from '@sapper/app';

	let user = null;
	let tname = "";
	let tdescr = "";
	let tcode = "";

    onMount(()=>{
        user = getUser();
    })

	async function create_test() {
		const res = await post(user, 'create_test', {
				title : tname,
				descr : tdescr,
				code  : tcode,
			})
		if (res.status == 200) {
			goto('/')
		}
	}

</script>

<style>
</style>

<title>New test</title>

<h4>Title</h4>
<input bind:value={tname} placeholder="Test Title">

<h4>Description</h4>
<textarea style="width:100%; height:200px" bind:value={tdescr} />

<h4>Code</h4>
<textarea style="width:100%; height:200px" bind:value={tcode} />

<button on:click={()=>create_test()} disabled={tname.length == 0 || tdescr.length == 0 || tcode.length == 0}>
	add test
</button>

