<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
body { 
  margin: 0;
}
#app {
  font-family: Roboto, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  margin: 0px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}
.router-view {
  min-height: 100vh;
}
</style>
