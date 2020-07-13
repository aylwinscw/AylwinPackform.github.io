<template>
  <b-container fluid>
    <b-row>
      <b-col col md="12">
        <h1>Orders</h1>
      </b-col>
    </b-row>
    <b-row align-v="center">
        <b-col col md="3">
            <h4><b-icon-search></b-icon-search>Search</h4>
            
        </b-col>
        <b-col col md="9">
            <b-form-input
              v-model="filter"
              type="search"
              id="filterInput"
              placeholder="Type to Search"
            ></b-form-input>
        </b-col>
    </b-row>
    <b-row>
        <b-col col md="3">
            <p><b>Created Date</b></p>
        </b-col>
    </b-row>
    <b-row>
        <b-col col md="12">
            <p><b>Total Amount</b>     ${{ total_amount }}</p>
        </b-col>
    </b-row>
    <b-row>
        <b-col col md="12">
            <b-table :items="orders"
            :current-page="currentPage"
            :per-page="perPage"
            :filter="filter"></b-table>
        </b-col>
    </b-row>
    <b-row>
        <b-col col md="12">
            <!-- <p>Total 2</p> -->
            <b-pagination
          v-model="currentPage"
          :total-rows="totalRows"
          :per-page="perPage"
          align="fill"
          size="sm"
          class="my-0"
        ></b-pagination>
        </b-col>
    </b-row>
  </b-container>
</template>

<script>
import Vue from 'vue'
import { BootstrapVue, BootstrapVueIcons } from 'bootstrap-vue'
import axios from 'axios'

Vue.use(BootstrapVue)
Vue.use(BootstrapVueIcons)

export default {
    name: 'Orders',
    data() {
        return {
            fields: ['order_id',
            'order_name', 
            'customer_company', 
            'customer_name', 
            'order_date', 
            'delivered_amount', 
            'total_amount'],
            orders: [],
            totalRows: null,
            currentPage: 1,
            perPage: 5,
            filter: null,
            total_amount: 0,
            start_date: null,
            end_date: null
        }
    },
    created() {

        // Access within promise
        var that = this;
        var tempData = [];
        var tempOrders = [];
        var tempCustomers = [];
        var tempCompanies = [];
        var amount = 0;

        var promise1 = axios.get('https://radiant-falls-27028.herokuapp.com/orders')
        .then(function(res) {
            tempOrders = res.data;
        })
        .catch(err => console.log(err));

        var promise2 = axios.get('https://radiant-falls-27028.herokuapp.com/customers')
        .then(function(res) {
            tempCustomers = res.data;
        })
        .catch(err => console.log(err));
        var promise3 = axios.get('https://radiant-falls-27028.herokuapp.com/companies')
        .then(function(res) {
            tempCompanies = res.data;

        })
        .catch(err => console.log(err));

        Promise.all([promise1, promise2, promise3]).then(function() {

            for(var cursor of tempOrders)
            {
                let ref_company_id;

                let order = {
                    order_name: "",
                    customer_company: "",
                    customer_name: "",
                    order_date: "",
                    delivered_amount: 99.11,
                    total_amount: 99.11
                };

                for(var customer of tempCustomers)
                {
                    if(cursor.customer_id === customer.user_id)
                    {
                        order.customer_name = customer.name;
                        ref_company_id = customer.company_id;
                    }
                }

                for(var company of tempCompanies)
                {
                    if(company.company_id === ref_company_id)
                    {
                        order.customer_company = company.company_name;
                    }
                }
                order.order_name = cursor.order_name;
                order.order_date = new Intl.DateTimeFormat('en-US', { year: 'numeric', month: 'short', day: '2-digit' }).format(new Date(Date.parse(cursor.created_at)));
                
                tempData.push(order);
                amount += order.total_amount;

            }

            that.total_amount = amount;
            that.orders = tempData
            that.totalRows = tempData.length;
            console.log(tempData.length);

        });
    }
}
</script>

<style scoped>

</style>