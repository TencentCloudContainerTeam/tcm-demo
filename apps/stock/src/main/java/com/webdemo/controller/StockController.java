package com.webdemo.controller;

import com.google.common.collect.Maps;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.math.BigDecimal;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

@RestController
public class StockController {

    @GetMapping("/stock")
    public Map<Integer, Integer> sales(@RequestParam("ids") String ids) {
        Map<Integer, Integer> mockDbMap = Maps.newHashMap();
        mockDbMap.put(1, 1001);
        mockDbMap.put(2, 1002);
        mockDbMap.put(3, 1003);
        mockDbMap.put(4, 1004);
        mockDbMap.put(5, 1005);
        mockDbMap.put(6, 1006);
        mockDbMap.put(7, 1007);
        mockDbMap.put(8, 1008);
        mockDbMap.put(9, 1009);
        mockDbMap.put(10, 1010);
        mockDbMap.put(11, 1011);
        mockDbMap.put(12, 1012);
        mockDbMap.put(13, 1013);
        mockDbMap.put(14, 1014);
        mockDbMap.put(14, 1015);

        Map<Integer, Integer> resultMap = Maps.newHashMap();

        for (String id : ids.split(",")) {
            Integer value = mockDbMap.get(Integer.valueOf(id));
            if (value != null) {
                // BigDecimal multiply = new BigDecimal(value).multiply(new BigDecimal(10));
                // resultMap.put(Integer.valueOf(id), multiply.toString());
                resultMap.put(Integer.valueOf(id), value);
            }
        }

        return resultMap;
    }

}
