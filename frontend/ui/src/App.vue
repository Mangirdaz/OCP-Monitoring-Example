<template>
  <div id="app">
    <img :src="pic_url" />
    <headerbar></headerbar>
    <notes></notes>
    
  </div>
</template>

<script>
import headerbar from './components/headerbar'
import notes from './components/notes'

export default {
  name: 'app',
  components: {
    headerbar,
    notes
  },
    data() {
      return {
        quote: '',
        pic_url: process.env.IMG_URL
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
