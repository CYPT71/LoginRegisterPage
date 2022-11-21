<template>
  <div>
    <div>
      <h1>Hello {{ user.Username }}</h1>
      <div>
        <p>email</p>
        <input v-model="email" />
      </div>
      <button @click="update">Update</button>
    </div>
    <div>
      <div class="exercices">
        <div class="python">
          <p>
            Lorem ipsum dolor sit amet consectetur adipisicing elit. Natus
            suscipit recusandae accusantium animi facere fugit. Ducimus nobis
            exercitationem corporis laudantium cupiditate repellendus
            consequuntur laborum aliquam iste minus vel ea, perspiciatis eos
            blanditiis iusto deleniti natus veniam in at, totam dolores?
          </p>
          <input type="file" @change="uploadFile" ref="file" />
          <button>refresh/submit</button>
        </div>
        <div class="python">
          <p>
            Lorem ipsum dolor sit amet consectetur adipisicing elit. Natus
            suscipit recusandae accusantium animi facere fugit. Ducimus nobis
            exercitationem corporis laudantium cupiditate repellendus
            consequuntur laborum aliquam iste minus vel ea, perspiciatis eos
            blanditiis iusto deleniti natus veniam in at, totam dolores?
          </p>
          <input type="file" @change="uploadFile" ref="file" />
          <button>refresh/submit</button>
        </div>
        <div class="c">
          <p>
            Lorem ipsum dolor sit amet consectetur adipisicing elit. Natus
            suscipit recusandae accusantium animi facere fugit. Ducimus nobis
            exercitationem corporis laudantium cupiditate repellendus
            consequuntur laborum aliquam iste minus vel ea, perspiciatis eos
            blanditiis iusto deleniti natus veniam in at, totam dolores?
          </p>
          <input type="file" @change="uploadFile" ref="file" />
          <button>refresh/submit</button>
        </div>
        <div class="c">
          <p>
            Lorem ipsum dolor sit amet consectetur adipisicing elit. Natus
            suscipit recusandae accusantium animi facere fugit. Ducimus nobis
            exercitationem corporis laudantium cupiditate repellendus
            consequuntur laborum aliquam iste minus vel ea, perspiciatis eos
            blanditiis iusto deleniti natus veniam in at, totam dolores?
          </p>

          <input type="file" @change="uploadFile" ref="file"/>
          <button>refresh/submit</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "@vue/reactivity";
import axios from "axios";

var user = ref({});
var email = ref();
const file = ref(null)

const props = defineProps({
  user: Object,
});

const uploadFile = (event) => {
  submitFile(event.target.files[0])
};

const submitFile = async (file) => {
  const formData = new FormData();
  formData.append("file", file);
  const headers = { "Content-Type": "multipart/form-data" };
  await axios.post("https://httpbin.org/post", formData, { headers }).then((res) => {
    res.data.files; // binary representation of the file
    res.status; // HTTP status
  });
};

const update = async () => {
  if (email.value) user.value.Email = email.value;

  axios.patch(
    "http://localhost:3000/user",
    {
      password: "",
      email: email.value,
    },
    {
      headers: { Authorization: "Bearer " + sessionStorage.token },
    }
  );
};
</script>

<style scoped lang="sass">
.exercices
  display: flex
  gap: 3rem
  margin-top: 2rem
.c .python
  gap: 2rem
</style>