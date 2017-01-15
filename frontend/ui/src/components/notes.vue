<template>
  <div class="notes">
    <h1>{{ msg }}</h1>
     <button class="btn btn-primary" v-on:click="getNotes()">Get Notes</button>
    <br />
    <br />
    Create New Note and hit enter
    <form @submit.prevent="createNote">
      <input type='text' v-model="newNote.note">
    </form>
     <br />
    <br />

    <table class="table table-bordered table-striped">
          <thead>
            <tr>
             <th class="col-sm-4">Note</th>
             <th class="col-sm-4">Creation Date</th>
             <th class="col-sm-2">ID</th>
             <th class="col-sm-2">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in notes">
              {{item}}
               <td>{{ item.note }}</td> 
                <td>{{ item.creation_time }}</td> 
                <td>{{ item.key }}</td> 
                <td>Delete</td> 
            </tr>
      </tbody>  
    </table>
   

  </div>
</template>

<script>
export default {
  name: 'notes',
  data: function(){
    return {
      msg: 'Notebook',
      newNote: {
            note: '',
        },
        notes: []
    }
  },
  methods: {
         createNote: function(e) {
               e.preventDefault();    

               this.$http.post(process.env.API_URL+"notes", this.newNote).then(function(response, status, request) {
                    console.log(response);
                }, function() {
                    console.log('failed');
                });
            },
        getNotes: function(){
        console.log("Getting notes "+ process.env.API_URL+"notes")
          this.$http.get(process.env.API_URL+"notes").then((response) => {
              // success callback
              var temp = JSON.parse(response.body)

              console.log(temp['Note'])
       
              var notes = [];
              temp['Note'].forEach(function(entry) {
                  notes.push(JSON.parse(entry.note))
              });
              this.notes= notes

            }, (response) => {
              // error callback
              console.log(response)
            });
      }
 }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}
</style>
